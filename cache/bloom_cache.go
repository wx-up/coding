package cache

import (
	"context"

	"github.com/wx-up/coding/cache/internal/errs"
)

// BloomCache 装饰器模式优化
type BloomCache struct {
	BloomFilter
	Cache
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (c *BloomCache) Get(ctx context.Context, key string) Value {
	val := c.Cache.Get(ctx, key)
	if val.Err != nil {
		return val
	}

	// 先查询 布隆过滤器
	exist := c.BloomFilter.Exist(key)
	if !exist {
		return Value{
			Err: errs.NewErrKeyNotFound(key),
		}
	}
	res, err := c.LoadFunc(ctx, key)
	return Value{
		Val: res,
		Err: err,
	}
}

type BloomFilter interface {
	Exist(string) bool
}
