package orm

// Aggregate 聚合函数
// SUM(age) COUNT(id) MIN(age) 等等
type Aggregate struct {
	// 可以是列名，也可以是复杂的表达式
	arg string
	// 聚合函数名称
	fn string
}

func (a Aggregate) selectable() {}

func Count(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "COUNT",
	}
}

func Avg(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "AVG",
	}
}
func Sum(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "SUM",
	}
}

func Min(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MIN",
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MAX",
	}
}
