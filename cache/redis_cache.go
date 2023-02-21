package cache

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache/internal/errs"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -destination=mocks/redis_cache.gen.go -package=mocks github.com/redis/go-redis/v9 Cmdable
type RedisCache struct {
	cmd redis.Cmdable
}

// NewRedisCacheV1 这是一种不好的设计，相当于在内部创建依赖，很难测试
// 如果用户非要使用这种方法，可以按照下面的实现，再委托给 NewRedisCache 方法
// 实际测试的时候，只测试 NewRedisCache 方法
//func NewRedisCacheV1(addr string) *RedisCache {
//	client := redis.NewClient(&redis.Options{
//		Addr: addr,
//	})
//	return NewRedisCache(client)
//}

// NewRedisCache 以 redis.Cmdable 作为参数的好处就是单元测试的时候我们可以 mock 测试 RedisCache
func NewRedisCache(cmd redis.Cmdable) *RedisCache {
	return &RedisCache{
		cmd: cmd,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) Value {
	res, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		return Value{
			Err: err,
		}
	}
	return Value{
		Val: res,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	res, err := r.cmd.Set(ctx, key, val, expiration).Result()
	if err != nil {
		return err
	}

	// 正常不需要判断 res 的值，这样子做了肯定万无一失
	if res != "OK" {
		return errs.ErrSetKeyFail
	}

	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.cmd.Del(ctx, key).Result()
	return err
}
