package orm

type Deletion[T any] struct {
	tbl string

	db *DB

	// WHERE 条件
	ps []Predicate
}

func NewDeletion[T any](db *DB) *Deletion[T] {
	return &Deletion[T]{
		db: db,
	}
}

func (d *Deletion[T]) Build() (*Query, error) {

	return nil, nil
}

func (d *Deletion[T]) From(tbl string) *Deletion[T] {
	d.tbl = tbl
	return d
}

func (d *Deletion[T]) Where(ps ...Predicate) *Deletion[T] {
	d.ps = ps
	return d
}
