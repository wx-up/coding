package cache

import (
	"context"
	"sync"
	"time"

	"github.com/wx-up/coding/cache/internal/errs"
)

type LocalCacheV1 struct {
	m         sync.Map
	close     chan struct{}
	closeOnce sync.Once
}

func (c *LocalCacheV1) Get(ctx context.Context, key string) Value {
	res, ok := c.m.Load(key)
	if !ok {
		return Value{
			Err: errs.NewErrKeyNotFound(key),
		}
	}
	return Value{
		Val: res,
	}
}

func (c *LocalCacheV1) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.m.Store(key, val)
	return nil
}

func (c *LocalCacheV1) Delete(ctx context.Context, key string) error {
	c.m.Delete(key)
	return nil
}

func (c *LocalCacheV1) Close() error {
	c.closeOnce.Do(func() {
		c.close <- struct{}{}
		close(c.close)
	})
	return nil
}
