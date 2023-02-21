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
	Close() error
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
