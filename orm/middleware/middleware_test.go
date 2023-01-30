package middleware

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/wx-up/coding/orm"
	"testing"
	"time"
)

// TestLog 测试 SQL 的打印
func TestLog(t *testing.T) {
	type TestModel struct {
	}
	logMiddleware := NewLogMiddlewareBuilder(func(query string, args ...any) {
		// 这里可以使用日志框架打印日志
		fmt.Println(query)
	}).Build()
	db, err := orm.OpenDB(nil, orm.DBWithMiddlewares([]orm.Middleware{logMiddleware}))
	require.NoError(t, err)
	_, err = orm.NewSelector[TestModel](db).Get(context.Background())
	fmt.Println(err)
}

// TestSlowQuery 测试慢 SQL 的打印
func TestSlowQuery(t *testing.T) {
	type TestModel struct {
	}
	// 慢sql阈值为 100ms
	slowQueryLog := NewSlowQueryMiddlewareBuilder(100, func(query string, args ...any) {
		fmt.Println(query)
	}).Build()

	db, err := orm.OpenDB(nil, orm.DBWithMiddlewares([]orm.Middleware{slowQueryLog, func(handler orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// sleep 1s 模拟长SQL执行
			time.Sleep(time.Second)
			return handler(ctx, qc)
		}
	}}))
	require.NoError(t, err)
	_, err = orm.NewSelector[TestModel](db).Get(context.Background())
}
