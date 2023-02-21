package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/wx-up/coding/cache/mocks"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/redis/go-redis/v9"
)

func Test(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	fmt.Println(n)
}

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name       string
		mock       func() redis.Cmdable
		key        string
		val        any
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "return OK",
			mock: func() redis.Cmdable {
				client := mocks.NewMockCmdable(ctrl)
				res := redis.NewStatusCmd(nil)
				res.SetVal("OK")
				res.SetErr(nil)
				client.EXPECT().Set(gomock.Any(), "key1", "val1", time.Second).Return(res)
				return client
			},
			key:        "key1",
			val:        "val1",
			expiration: time.Second,
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewRedisCache(tc.mock())
			err := cache.Set(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
