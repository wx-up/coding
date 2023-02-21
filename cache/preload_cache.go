package cache

import (
	"context"
	"fmt"
	"time"
)

// PreloadCache 某个 key 快过期了，则重新加载数据，并设置缓存
// 借助哨兵模式
type PreloadCache struct {
	cache         Cache
	sentinelCache *LocalCache

	LoadFunc func(ctx context.Context, key string) (any, error)
}

func NewPreloadCache(cache Cache) *PreloadCache {
	res := &PreloadCache{
		cache: cache,
	}

	// 哨兵 cache
	// 当 key 过期的时候，从数据库中重新捞数据，并设置缓存
	res.sentinelCache = NewLocalCache(WithOnEvicted(func(key string, val any) {
		data, err := res.LoadFunc(context.Background(), key)
		if err != nil {
			fmt.Println(err)
		}
		err = res.Set(context.Background(), key, data, time.Second)
		if err != nil {
			fmt.Println(err)
		}
	}))
	return res
}

func (c *PreloadCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 提前5秒过期
	_ = c.sentinelCache.Set(ctx, key, "", expiration-time.Second*5)

	return c.cache.Set(ctx, key, val, expiration)
}
