package pattern

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache"
)

type Logger interface {
	Error()
}

var defaultLogger Logger

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// ReadThroughCache read-through 模式
// 读操作：从缓存中读取，读取不到则从 "DB" 中捞数据，并设置缓存，再返回
// 写操作：业务代码自己负责写缓存和写数据库
type ReadThroughCache struct {
	cache.Cache

	// 将 "从数据库捞数据" 的动作抽象为一个方法，实际是从文件、数据库获取数据我们并不关心
	Loader

	// 过期时间
	Expiration time.Duration

	logger Logger
}

func (c *ReadThroughCache) Get(ctx context.Context, key string) cache.Value {
	val := c.Cache.Get(ctx, key)
	if val.Err == nil {
		return val
	}
	// 从数据库捞
	data, err := c.Loader.Load(ctx, key)
	if err != nil {
		// 1. 如果你的中间件/类有日志抽象，那么可以打印操作，然后返回一个新的错误
		// c.log.Println(err)
		// return nil, errors.New("无法加载数据")
		// 2. 否则，你就不应该丢掉原始错误信息（ 一般采用 wrap 的方式）

		// 错误最好包装一下
		// fmt.Errorf("read-through: 无法加载数据 %w",err)
		// 这里不能直接 errors.New 一个新的错误，这就相当于把 根本错误 丢弃了，后续不好调试
		return cache.Value{
			Err: err,
		}
	}
	res := cache.Value{
		Val: data,
	}
	// 设置缓存
	// 出错了：
	//  1. 可以选择打印一条 warn 日志
	//  2. 忽略，并累计错误，达到阈值的时候，报警
	_ = c.Set(ctx, key, res, c.Expiration)

	return res
}

type Loader interface {
	Load(ctx context.Context, key string) (any, error)
}

type LoadFunc func(ctx context.Context, key string) (any, error)

func (f LoadFunc) Load(ctx context.Context, key string) (any, error) {
	return f(ctx, key)
}
