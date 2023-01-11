package orm

// Assignable 标记接口，用于 Insert 语句
type Assignable interface {
	assign()
}

type Assigment struct {
	column string
	val    any
}

func (a Assigment) assign() {}

func Assign(column string, val any) Assigment {
	return Assigment{
		column: column,
		val:    val,
	}
}
