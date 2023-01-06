package orm

type Model struct {
	// 结构体对应的表名
	tableName string
	// 结构体所有字段的元数据信息
	fieldMap map[string]*field
}

type field struct {
	// 当前字段对应的列名
	colName string
}
