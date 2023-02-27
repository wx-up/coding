package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/wx-up/coding/cache/internal/errs"
)

/*
内存控制的两种策略：
	1. 控制键值对的数量
	2. 控制整体内存

装饰器实现

*/

// LocalCacheMaxCnt 控制键值对数量的实现
type LocalCacheMaxCnt struct {
	maxCnt int32
	cnt    int32
	*LocalCache
}

func NewLocalCacheMaxCnt(localCache *LocalCache) *LocalCacheMaxCnt {
	res := &LocalCacheMaxCnt{
		maxCnt: 100,
	}
	origin := localCache.OnEvicted
	localCache.OnEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	res.LocalCache = localCache
	return res
}

// Set TODO：重复设置某一个 key 的时候， cnt 不应该增加
func (c *LocalCacheMaxCnt) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 先占一个位置
	cnt := atomic.AddInt32(&c.cnt, 1)
	if cnt > c.maxCnt {
		// 如果已经满了，就把1减回去
		atomic.AddInt32(&c.cnt, -1)
		return errs.ErrCacheIsCompletely
	}
	return c.LocalCache.Set(ctx, key, val, expiration)
}

//func (c *LocalCacheMaxCnt) Delete(ctx context.Context, key string) error {
//	// 这种实现存在问题：
//	// 因为在 Cache 的实现中 Delete 方法不是删除的唯一入口，Get 的时候以及 goroutine 轮训的时候都有可能删除 key
//	// 解决方式：使用 CDC 机制，onEvicted 方法
//	defer func() {
//		atomic.AddInt32(&c.cnt, -1)
//	}()
//	return c.Cache.Delete(ctx, key)
//}
