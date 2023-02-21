package v1

import (
	"context"
	"time"

	"github.com/wx-up/coding/cache/lock/redis/internal/errs"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client redis.Cmdable
}

func (c *Client) TryLock(ctx context.Context, key string, val any, expiration time.Duration) error {
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errs.ErrFailedToPreemptLock
	}
	return nil
}

func (c *Client) UnLock(ctx context.Context, key string) error {
	res, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	// 删除成功返回 1
	// 删除失败返回 0 比如 key 已经过期
	if res != 1 {
		return errs.ErrLockNotHand
	}
	return nil
}
