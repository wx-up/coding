package orm

import (
	"context"
)

// Querier select 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor delete、update、insert 语句
type Executor interface {
	Exec(ctx context.Context) Result
}

// QueryBuilder 构建 SQL
type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
