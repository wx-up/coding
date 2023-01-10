package orm

// expression 标记接口
// expr 方法并不具有实际含义和实际用途，仅仅是标记
type expression interface {
	expr()
}

type RawExpr struct {
	raw  string
	args []any
}

func Raw(raw string, args ...any) RawExpr {
	return RawExpr{
		raw:  raw,
		args: args,
	}
}

func (r RawExpr) expr() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

func (r RawExpr) selectable() {}
