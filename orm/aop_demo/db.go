package aop_demo

import (
	"context"
	"database/sql"
)

/*
 如果第三方库没有支持 AOP 方案，我们可以通过装饰器的方式实现，比如 sql.DB，可以参考 web 的 middleware 实现
 AOP：说白了就是前面执行一些什么、后面执行一些什么
*/

// AopDB 将 sql.DB 结构体所有的公开方法实现一遍（ 内部委托给db 对象 ）
type AopDB struct {
	db *sql.DB
	ms []Middleware
}

type AopDBContext struct {
	query string
	args  []any
}

type Handler func(ctx *AopDBContext) AopDBResult

type Middleware func(next Handler) Handler

type AopDBResult struct {
	result sql.Result
	err    error
}

func (db *AopDB) ExecContext(c context.Context, query string, args ...any) (sql.Result, error) {
	handler := func(ctx *AopDBContext) (res AopDBResult) {
		res.result, res.err = db.db.ExecContext(c, ctx.query, ctx.args...)
		return
	}

	for i := len(db.ms); i >= 0; i-- {
		handler = db.ms[i](handler)
	}

	res := handler(&AopDBContext{
		query: query,
		args:  args,
	})
	return res.result, res.err
}
