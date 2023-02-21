package cache

import (
	"context"
	"sync"
	"time"
)

type LocalCache struct {
	lock      sync.RWMutex
	m         map[string]*item
	close     chan struct{}
	closeOnce sync.Once
	OnEvicted func(key string, val any)
}

type LocalCacheOption func(*LocalCache)

func WithOnEvicted(f func(key string, val any)) LocalCacheOption {
	return func(cache *LocalCache) {
		cache.OnEvicted = f
	}
}

func NewLocalCache(opts ...LocalCacheOption) *LocalCache {
	res := &LocalCache{
		m:     make(map[string]*item),
		close: make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (c *LocalCache) Get(ctx context.Context, key string) Value {
	// TODO implement me
	panic("implement me")
}

func (c *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (c *LocalCache) Delete(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (c *LocalCache) delete(key string, val any) {
	delete(c.m, key)
	if c.OnEvicted != nil {
		c.OnEvicted(key, val)
	}
}
