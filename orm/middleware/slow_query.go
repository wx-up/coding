package middleware

import (
	"context"
	"github.com/wx-up/coding/orm"
	"time"
)

// SlowQueryMiddlewareBuilder slow query log
type SlowQueryMiddlewareBuilder struct {
	// 阈值，毫秒
	threshold int64

	logFunc func(query string, args ...any)
}

func NewSlowQueryMiddlewareBuilder(threshold int64, f func(query string, args ...any)) *SlowQueryMiddlewareBuilder {
	return &SlowQueryMiddlewareBuilder{
		threshold: threshold,
		logFunc:   f,
	}
}

func (b *SlowQueryMiddlewareBuilder) Build() orm.Middleware {
	return func(handler orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			start := time.Now()
			query, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}

			// 当超过阈值的时候，记录日志
			defer func() {
				cost := time.Now().Sub(start)
				if b.threshold > 0 && cost.Milliseconds() >= b.threshold {
					b.logFunc(query.SQL, query.Args...)
				}
			}()

			return handler(ctx, qc)
		}
	}
}
