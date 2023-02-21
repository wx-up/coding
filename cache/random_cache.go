package cache

import (
	"context"
	"math/rand"
	"time"
)

type RandomCache struct {
	Cache
}

func (c *RandomCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	offset := rand.Intn(300)
	expiration = expiration + time.Duration(offset)*time.Second
	return c.Cache.Set(ctx, key, val, expiration)
}
