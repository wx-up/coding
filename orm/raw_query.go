package orm

import "context"

// RawQuerier 原生查询
// 用户直接使用原生sql或者通过Selector等构建SQL之后，最后都会通过 RawQuerier 和数据库交互，返回结果集
type RawQuerier[T any] struct {
	sql  string
	args []any

	sess Session

	ms []Middleware
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{SQL: r.sql, Args: r.args}, nil
}

func RawQuery[T any](sess Session, sql string, args ...any) *RawQuerier[T] {
	core := sess.getCore()
	return &RawQuerier[T]{
		sql:  sql,
		args: args,
		sess: sess,
		ms:   core.Middlewares(),
	}
}

// Use 添加 RawQuerier 层面的 middleware
func (r *RawQuerier[T]) Use(ms []Middleware) *RawQuerier[T] {
	r.ms = append(r.ms, ms...)
	return r
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	return r.get(ctx)
}

func (r *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("imp")
}

func (r *RawQuerier[T]) Exec(ctx context.Context) Result {
	return Result{}
}

func (r *RawQuerier[T]) get(ctx context.Context) (*T, error) {
	root := r.getHandler(r.sess)
	for i := len(r.ms) - 1; i >= 0; i-- {
		root = r.ms[i](root)
	}

	res := root(ctx, &QueryContext{
		Type:    QueryTypeSelect,
		Builder: r, // r 实现了 QueryBuilder 接口
	})

	if res.Err != nil {
		return nil, res.Err
	}

	if res.Result != nil {
		// 这里可能会 panic 例如用户自定义中间件时，给 Result 赋值为非 T 类型时
		// 但是这是用户使用不当导致的，框架设计者可以不处理
		return res.Result.(*T), nil
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) getHandler(sess Session) Handler {
	return func(ctx context.Context, qc *QueryContext) *QueryResult {
		query, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		rows, err := sess.queryContext(ctx, query.SQL, query.Args...)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}

		if !rows.Next() {
			return &QueryResult{
				Err: ErrNoRows,
			}
		}

		t := new(T)
		model, err := sess.getCore().R().Get(t)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		err = sess.getCore().ValCreator()(model, t).SetColumns(rows)

		return &QueryResult{
			Result: t,
			Err:    err,
		}
	}
}
