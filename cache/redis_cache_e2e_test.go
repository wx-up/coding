//go:build e2e

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/redis/go-redis/v9"
)

// TestRedisCache_e2e_Set 集成测试
func TestRedisCache_e2e_Set(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	cache := NewRedisCache(client)
	err := cache.Set(context.Background(), "key", "value", time.Second)
	require.NoError(t, err)

	val := cache.Get(context.Background(), "key")
	require.NoError(t, val.Err)
	require.Equal(t, "value", val.Val)
}
