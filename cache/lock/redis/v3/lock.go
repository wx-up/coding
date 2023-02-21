package v3

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/wx-up/coding/cache/lock/redis/internal/errs"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

var (
	//go:embed script/unlock.lua
	luaUnLock string
	//go:embed script/refresh.lua
	refresh string
	//go:embed script/lock.lua
	lock string
)

/*
 lua 脚本解决并发问题：
	lua 脚本将多个操作包装到一起，因为 redis 是单线程执行的，所以封装的操作之间一定是时序执行的，中间插入不了其他操作
	如果其他缓存是多线程的，那将多个操作封装是没有什么用
*/

type Client struct {
	cmd redis.Cmdable

	g singleflight.Group
}

func NewClient(cmd redis.Cmdable) *Client {
	return &Client{
		cmd: cmd,
	}
}

func (c *Client) SingleFlightLock(ctx context.Context, key string, expiration time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	for {
		flag := false

		result := c.g.DoChan(key, func() (interface{}, error) {
			// 当多个 goroutine 准备抢锁时，只有一个 goroutine 会进入执行函数体
			flag = true
			// 请求 redis 设置 key - value
			return c.Lock(ctx, key, expiration, retry, timeout)
		})
		select {
		case res := <-result:
			// 只有 flag 标记为 true 的 goroutine 才是真正抢到锁的，并返回
			// 其他 flag 标记为 false 的，继续下一轮抢锁
			if flag {
				// 这里最好根据具体的错误走不同的逻辑
				if res.Err != nil {
					return nil, res.Err
				}
				c.g.Forget(key)
				return res.Val.(*Lock), nil
			}
		case <-ctx.Done(): // 超时就直接退出
			return nil, ctx.Err()
		}
	}
}

// Do 这种设计有点类似 MYSQL 的事务 API
func (c *Client) Do(ctx context.Context, f func()) error {
	// 随机生成一个 key
	key := "rand_key"
	expiration := time.Second * 2
	l, err := c.TryLock(ctx, key, expiration)
	if err != nil {
		return err
	}

	// 自动续约
	go func() {
		err := l.AuthRefresh(expiration-time.Second, expiration-time.Second)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 执行业务代码
	f()

	return l.UnLock(ctx)
}

// Lock 是尽可能重试减少加锁失败的可能
// Lock 会在超时或者锁正被人持有的时候进行重试
// 最后返回的 error 使用 errors.Is 判断，可能是：
// - context.DeadlineExceeded: Lock 整体调用超时
// - ErrFailedToPreemptLock: 超过重试次数，但是整个重试过程都没有出现错误
// - DeadlineExceeded 和 ErrFailedToPreemptLock: 超过重试次数，但是最后一次重试超时了
// 你在使用的过程中，应该注意：
// - 如果 errors.Is(err, context.DeadlineExceeded) 那么最终有没有加锁成功，谁也不知道
// - 如果 errors.Is(err, ErrFailedToPreemptLock) 说明肯定没成功，而且超过了重试次数
// - 否则，和 Redis 通信出了问题
func (c *Client) Lock(ctx context.Context, key string,
	expiration time.Duration, retryStrategy RetryStrategy, timeout time.Duration,
) (*Lock, error) {
	val := uuid.New().String()

	// 刚进来的时候，检测以下 ctx
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		ctx1, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.cmd.Eval(ctx1, lock, []string{key}, val, expiration.Seconds()).Result()
		cancel()

		// 成功获得锁
		if res == "OK" {
			return &Lock{
				cmd:        c.cmd,
				val:        val,
				key:        key,
				expiration: expiration,
				unlock:     make(chan struct{}, 1),
			}, nil
		}

		// 非超时错误，那么基本上代表遇到了一些不可挽回的场景，所以没太大必要继续尝试了
		// 比如说 Redis server 崩了，或者 EOF 了
		// 最后使用 errors.Is 判断错误，如果使用 == 的话，如果别人 wrap 过的话，就会返回 false
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		interval, ok := retryStrategy.Next()
		if !ok {
			if err != nil {
				err = fmt.Errorf("最后一次重试错误：%w", err)
			} else {
				err = fmt.Errorf("锁被人持有：%w", err)
			}
			return nil, fmt.Errorf("redis-lock：重试机会耗尽，%w", err)
		}

		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}

		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	ok, err := c.cmd.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errs.ErrFailedToPreemptLock
	}
	return &Lock{
		cmd:        c.cmd,
		val:        val,
		key:        key,
		expiration: expiration,
	}, nil
}

type Lock struct {
	cmd        redis.Cmdable
	val        string
	key        string
	expiration time.Duration
	unlock     chan struct{}
	unlockOnce sync.Once
}

// AuthRefresh 自动续约
// internal 每隔多久续约一次，也可以根据 expiration 计算，比如它的1/2或者2/3
// timeout redis 操作的超时时间
func (l *Lock) AuthRefresh(internal time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(internal)
	// retrySignal 重试信号
	retrySignal := make(chan struct{}, 1)
	defer func() {
		ticker.Stop()
		close(retrySignal) // 不关其实也没有关系，最后反正会被垃圾回收器回收掉
	}()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()

			// 超时了，就重试
			if err == context.DeadlineExceeded {
				retrySignal <- struct{}{}
				continue
			} else {
				return err
			}
		case <-retrySignal:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()

			// 超时了，就重试
			if err == context.DeadlineExceeded {
				retrySignal <- struct{}{}
				continue
			} else {
				return err
			}
		case <-l.unlock:
			return nil

		}
	}
}

// Refresh 手动续约
func (l *Lock) Refresh(ctx context.Context) error {
	// 过期时间转成秒，一般业务场景过期时间都是以秒为单位
	res, err := l.cmd.Eval(ctx, refresh, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return errs.ErrLockNotHand
	}
	return nil
}

func (l *Lock) UnLock(ctx context.Context) error {
	// 防止用户多次调用 Unlock，虽然我们期望用户只会调用一次
	l.unlockOnce.Do(func() {
		// l.unlock <- struct{}{}
		// 如果 AuthRefresh 被调用一次（ 一般在 goroutine 中调用 ），这里发送一个信号是没有问题的
		// 但是如果用户调用 AuthRefresh 很多次，而这里只发送一个信号，就会导致其他 AuthRefresh 不会退出
		// close 的话，所有监听的 chan 都会收到消息，所以不存在上述的问题
		close(l.unlock)
	})
	// 发送了 luaUnLock 脚本到 redis，redis 再执行脚本，注意不是 Go 调度执行的
	res, err := l.cmd.Eval(ctx, luaUnLock, []string{l.key}, l.val).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return errs.ErrLockNotHand
	}
	return nil
}
