package cache

import (
	"context"
	"errors"
	"time"
)

//go:generate mockgen -destination=mocks/cache.gen.go -package=mocks -source=cache.go Cache
type Cache interface {
	Get(ctx context.Context, key string) Value
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type Value struct {
	Val any
	Err error
}

func (v Value) String() (string, error) {
	if v.Err != nil {
		return "", v.Err
	}
	res, ok := v.Val.(string)
	if !ok {
		return "", errors.New("不可转换类型")
	}
	return res, nil
}

type CacheV2 interface {
	Get(ctx context.Context, key string) Value
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	OnEvicted(func(key string, val any))
}

type CacheV3 interface {
	Get(ctx context.Context, key string) Value
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error

	// Subscribe channel 的方案有一点不是很好，就是缓冲区的大小不好设置
	// 如果设置的比较小，消费端又处理的比较慢，就会阻塞正常的缓存操作
	Subscribe() <-chan Event
}

type Event struct {
	Key  string
	Val  any
	Type int8
}
