package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrCacheClosed = errors.New("cache：缓存已经关闭")

// CloseCache 装饰器模式，实现一个可以关闭的 cache
type CloseCache struct {
	Cache
	closed bool
	lock   sync.RWMutex
}

func (c *CloseCache) Get(ctx context.Context, key string) Value {
	c.lock.RLock()
	if c.closed {
		return Value{
			Err: ErrCacheClosed,
		}
	}
	c.lock.RUnlock()
	return c.Cache.Get(ctx, key)
}

func (c *CloseCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.lock.RLock()
	if c.closed {
		return ErrCacheClosed
	}
	c.lock.RUnlock()
	return c.Cache.Set(ctx, key, val, expiration)
}

func (c *CloseCache) Delete(ctx context.Context, key string) error {
	c.lock.RLock()
	if c.closed {
		return ErrCacheClosed
	}
	c.lock.RUnlock()
	return c.Cache.Delete(ctx, key)
}

func (c *CloseCache) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closed = true
}
