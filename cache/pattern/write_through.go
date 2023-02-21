package pattern

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache"
)

// WriteThroughCache 目前的实现都是同步的
// 异步、半异步实现可以采用标记位的方式，比如 Async bool 但是标记位写的代码比较丑
// 采用装饰器的方式，但是比较难实现，特别是 StoreFunc 的处理
type WriteThroughCache struct {
	cache.Cache
	StoreFunc
}

// Set 先写缓存再写数据库，还是先写数据库再写缓存 其实都可以的
// 数据不一致的问题始终存在
func (c *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := c.Cache.Set(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	if c.StoreFunc != nil {
		err = c.StoreFunc(ctx, key, val)
	}
	return err
}

type Store interface {
	Store(ctx context.Context, key string, val any) error
}

type StoreFunc func(ctx context.Context, key string, val any) error

func (f StoreFunc) Store(ctx context.Context, key string, val any) error {
	return f(ctx, key, val)
}
