package middleware

import (
	"context"

	"github.com/wx-up/coding/orm"
)

// LogMiddlewareBuilder access log
type LogMiddlewareBuilder struct {
	logFunc func(query string, args ...any)
}

func NewLogMiddlewareBuilder(f func(query string, args ...any)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFunc: f, // f 内部的实现可以使用日志框架进行打印或者直接使用 fmt.Println
	}
}

func (lb LogMiddlewareBuilder) Build() orm.Middleware {
	return func(handler orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			b, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			lb.logFunc(b.SQL, b.Args...)
			return handler(ctx, qc)
		}
	}
}
