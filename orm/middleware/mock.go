package middleware

import (
	"context"
	"github.com/wx-up/coding/orm"
	"time"
)

type MockMiddlewareBuilder struct {
}

type Mock struct {
	Sleep  time.Duration
	Result *orm.QueryResult
}
type MockKey struct {
}

func (b *MockMiddlewareBuilder) Build() orm.Middleware {
	return func(handler orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 当 ctx 存在 mockKey 的时候，就会进入 mock 的逻辑
			val := ctx.Value(MockKey{})
			if val != nil {
				mock := val.(*Mock)
				if mock.Sleep > 0 {
					time.Sleep(mock.Sleep)
				}
				return mock.Result
			}
			return handler(ctx, qc)
		}
	}
}
