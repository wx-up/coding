package orm

type TableReference interface {
	tableAlias() string
}

// Table 普通表
type Table struct {
	entry any
	alias string
}

func (t Table) tableAlias() string {
	return t.alias
}

// TableOf  entry 为指针类型
func TableOf(entry any) Table {
	return Table{
		entry: entry,
	}
}

// As 不可变对象的设计，可以帮助编译器识别出对象不会逃逸到堆上，分配到栈上即可
func (t Table) As(alias string) Table {
	return Table{
		entry: t.entry,
		alias: alias,
	}
}

// Join 返回 builder（ 建造者模式往往是链式调用，返回值往往是指针类型 ）
func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "JOIN",
		right: right,
	}
}

// LeftJoin Join 是比较复杂的，你不能假定右边部分就一定是 Table，它可以是 subQuery 也可以是 Join
// 因此参数的类型应该是 TableReference
func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "LEFT JOIN",
		right: right,
	}
}

func (t Table) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "RIGHT JOIN",
		right: right,
	}
}

// Join 查询
type Join struct {
	left TableReference

	// JOIN、LEFT JOIN
	typ string

	right TableReference

	on    []Predicate
	using []string
}

// tableAlias join 语句是没有别名的
// (goods join user) a  这是不合法的
// 但是 join 内部的表是可以有别名的 (goods g join user u)
func (j Join) tableAlias() string {
	return ""
}

type JoinBuilder struct {
	left TableReference

	// JOIN、LEFT JOIN
	typ string

	right TableReference
}

func (jb *JoinBuilder) On(pc ...Predicate) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		on:    pc,
	}
}

func (jb *JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		using: cols,
	}
}

// SubQuery 子查询
type SubQuery struct {
}

func (s SubQuery) tableAlias() string {
	//TODO implement me
	panic("implement me")
}
