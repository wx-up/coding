package pattern

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache"
)

type CacheV2 interface {
	cache.Cache
	OnEvicted(func(key string, val any))
}
type WriteBackCache struct {
	CacheV2

	TimeOut time.Duration
}

func NewWriteBackCache(cache CacheV2, store func(ctx context.Context, key string, val any)) CacheV2 {
	res := &WriteBackCache{}
	cache.OnEvicted(func(key string, val any) {
		ctx, cancel := context.WithTimeout(context.Background(), res.TimeOut)
		defer cancel()
		store(ctx, key, val)
	})
	res.CacheV2 = cache

	return res
}

type WriteBackCacheV2 struct {
	*cache.LocalCache
	TimeOut time.Duration
}

func NewWriteBackCacheV2(localCache *cache.LocalCache, store func(ctx context.Context, key string, val any)) *WriteBackCacheV2 {
	res := &WriteBackCacheV2{
		LocalCache: localCache,
	}
	origin := localCache.OnEvicted
	localCache.OnEvicted = func(key string, val any) {
		if origin != nil {
			origin(key, val)
		}

		ctx, cancel := context.WithTimeout(context.Background(), res.TimeOut)
		defer cancel()
		store(ctx, key, val)
	}

	return res
}

func (c *WriteBackCacheV2) Close() error {
	// 遍历所有的 key 将值刷新到数据库中
	return nil
}
