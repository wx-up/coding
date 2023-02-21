package pattern

import (
	"context"
	"strings"
	"testing"
)

func TestReadThrough(t *testing.T) {
	cache := &ReadThroughCache{
		// 如果共用一个 read_through 的话，loadData 函数的实现会是一坨屎
		// 一般我们会选择一个类型的cache一个实例
		Loader: LoadFunc(func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 加载用户的数据
			} else if strings.HasPrefix(key, "/order/") {
				// 加载订单的数据
			}
			return nil, nil
		}),
	}
	_ = cache

	// 一个类型一个 cache
	userCache := &ReadThroughCache{
		Loader: LoadFunc(func(ctx context.Context, key string) (any, error) {
			if !strings.HasPrefix(key, "/user/") {
				return nil, nil
			}
			return nil, nil
		}),
	}
	_ = userCache
}
