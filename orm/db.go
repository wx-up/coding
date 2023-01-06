package orm

import "reflect"

type DB struct {
	r *Registry
}

type DBOption func(db *DB)

// NewDB 创建 DB 实例
// 预留了 error 公开方法尽量都加上 error
func NewDB(opts ...DBOption) (*DB, error) {
	db := &DB{
		r: &Registry{
			models: make(map[reflect.Type]*Model),
		},
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

func MustNewDB(opts ...DBOption) *DB {
	db, err := NewDB(opts...)
	if err != nil {
		panic(err)
	}
	return db
}
