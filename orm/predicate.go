package orm

type op string

const (
	opEq  op = "="
	opLt  op = "<"
	opLte op = "<="
	opGt  op = ">"
	opGte op = ">="
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (p op) String() string {
	return string(p)
}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{
		name: name,
	}
}

func (c Column) expr() {}

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
