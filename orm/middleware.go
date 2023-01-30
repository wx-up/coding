package orm

import "context"

const (
	QueryTypeSelect string = "SELECT"
	QueryTypeUpdate string = "UPDATE"
	QueryTypeInsert string = "INSERT"
)

type QueryContext struct {
	// 用来标记 SELECT、UPDATE、INSERT
	Type string

	// 通过断言 Selector、Inserter 等实现，来修改结构体内部的字段，从而达到修改 SQL 的目的
	Builder QueryBuilder
}

type QueryResult struct {
	// SELECT 语句，返回值为 T 或者 []T
	// UPDATE、DELETE、INSERT 语句返回值 Result
	Result any

	Err error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(handler Handler) Handler
