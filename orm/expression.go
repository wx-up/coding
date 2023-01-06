package orm

// expression 标记接口
// expr 方法并不具有实际含义和实际用途，仅仅是标记
type expression interface {
	expr()
}
