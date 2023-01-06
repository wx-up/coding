package predicate_v1

//type Predicate struct {
//	Column string
//	Op     string
//	Arg    any
//}
//
//// Eq 构造 Predicate 的方式一
//func Eq(column string, arg any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Arg:    arg,
//	}
//}

type op string

const (
	opEq  op = "="
	opLt  op = "<"
	opLte op = "<="
	opGt  op = ">"
	opGte op = ">="
	opNot op = "NOT"
	opAnd op = "AND"
)

type Predicate struct {
	column column
	op     op

	arg any
}

// Column 构造 Predicate 的方式二，将 column 包装成结构体，这里采用方式二
type column struct {
	name string
}

func C(name string) column {
	return column{name: name}
}

func (c column) Eq(arg any) Predicate {
	return Predicate{
		column: c,
		op:     opEq,
		arg:    arg,
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:  opNot,
		arg: p,
	}
}
