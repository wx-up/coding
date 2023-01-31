package orm

// RawQuery 原生查询
type RawQuery struct {
	sql  string
	args []any
}
