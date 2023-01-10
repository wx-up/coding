package orm

type op string

const (
	opEmpty op = ""
	opEq    op = "="
	opLt    op = "<"
	opLte   op = "<="
	opGt    op = ">"
	opGte   op = ">="
	opNot   op = "NOT"
	opAnd   op = "AND"
	opOr    op = "OR"
)

func (p op) String() string {
	return string(p)
}

type Column struct {
	name  string
	alias string
}

// C name 为字段名不是列名
// 对于 orm 来说，用户应该操作字段名，不必知道数据库的定义，从而达到解耦的效果
// 像 beego 和 gorm 其实操作的是列名
func C(name string) Column {
	return Column{
		name: name,
	}
}

// As 别名
func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

func (c Column) expr()       {}
func (c Column) selectable() {}

func (c Column) Eq(val any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		// 使用 Val 结构体包装一下，让它实现 expression 接口
		right: Val{val: val},
	}
}

func (c Column) Gt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: Val{val: val},
	}
}

func (c Column) Gte(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGte,
		right: Val{val: val},
	}
}

func (c Column) Lt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: Val{val: val},
	}
}

func (c Column) Lte(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLte,
		right: Val{val: val},
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

type Predicate struct {
	left  expression
	op    op
	right expression
}

func (p Predicate) expr() {}

func (p Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: right,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: right,
	}
}

type Val struct {
	val any
}

func (v Val) expr() {

}

type OrderBy struct {
	col   string
	order string
}

func Asc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "ASC",
	}
}

func Desc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "DESC",
	}
}
