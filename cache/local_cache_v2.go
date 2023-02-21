package cache

import (
	"context"
	"sync"
	"time"

	"github.com/wx-up/coding/cache/internal/errs"
)

type LocalCacheV2 struct {
	m         sync.Map
	close     chan struct{}
	closeOnce sync.Once
}

func NewLocalCacheV2() *LocalCacheV2 {
	c := &LocalCacheV2{
		close: make(chan struct{}, 1),
	}
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				cnt := 0
				c.m.Range(func(key, value any) bool {
					v := value.(*item)
					// time.Now 是一个比较慢的操作，这里会有性能瓶颈
					// 如果 time.Now 放在 Range 之前先求值，当 map 的元素比较多的时候就会存在误差
					if v.deadline.Before(time.Now()) {
						c.m.Delete(key)
					}
					cnt++
					return cnt < 1000
				})
			case <-c.close:
				ticker.Stop()
				return
			}
		}
	}()
	return c
}

func (c *LocalCacheV2) Get(ctx context.Context, key string) Value {
	val, ok := c.m.Load(key)
	if !ok {
		return Value{
			Err: errs.NewErrKeyNotFound(key),
		}
	}
	itm := val.(*item)
	if itm.deadline.Before(time.Now()) {
		c.m.Delete(key)
		return Value{
			Err: errs.NewErrKeyNotFound(key),
		}
	}
	return Value{
		Val: itm.val,
	}
}

func (c *LocalCacheV2) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (c *LocalCacheV2) Delete(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (c *LocalCacheV2) Close() error {
	c.closeOnce.Do(func() {
		c.close <- struct{}{}
		close(c.close)
	})
	return nil
}

// item 可以考虑使用 sync.Pool 来复用
type item struct {
	val      any
	deadline time.Time
}
