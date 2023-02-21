package v2

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache/lock/redis/internal/errs"

	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

/*
 当前版本的设计，没有考虑并发问题
*/

type Client struct {
	cmd redis.Cmdable
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
		cmd: c.cmd,
		val: val,
	}, nil
}

type Lock struct {
	cmd redis.Cmdable
	val string
}

func (l *Lock) UnLock(ctx context.Context, key string) error {
	val, err := l.cmd.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if l.val != val { // 代表锁不是你的锁
		return errs.ErrLockNotHand
	}

	// 删除锁
	res, err := l.cmd.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return errs.ErrLockNotHand
	}
	return nil
}
