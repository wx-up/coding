package orm

import (
	"database/sql"
	"github.com/wx-up/coding/orm/internal/model"
	"github.com/wx-up/coding/orm/internal/valuer"
)

// DB 它其实是 sql.DB 的一个装饰器
type DB struct {
	// 如果直接使用匿名组合的话，用户可以直接调用 sql.DB 的公开方法，从而绕过了 orm
	// 使用小写 db 是为了限制用户使用 orm 的方法操作数据库，而不是直接使用 sql.DB 的方法
	db *sql.DB
	r  *model.Registry

	valCreator valuer.Factory

	// 方言抽象应该放在 db 里面，因为它是属于 db 的，不同的 db 方言不同
	dialect Dialect
}

func (db *DB) Migrate(ms ...any) {

}

type DBOption func(db *DB)

// Open 创建 DB 实例
// 预留了 error 公开方法尽量都加上 error
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// OpenDB 直接传递一个 *sql.DB 来构建 DB
// 该方法的好处是方便 mock 测试
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:          model.NewRegister(),
		db:         db,
		valCreator: valuer.NewReflectValuer, // 函数式编程
		dialect:    NewMysqlDialect(),       // 指定方言，默认是 mysql
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithUnsafeValCreator() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewUnsafeValuer
	}
}

// DBWithDialect 指定方言
func DBWithDialect(d Dialect) DBOption {
	return func(db *DB) {
		db.dialect = d
	}
}

//func DBWithRegistry(r model.RegistryInterface) DBOption {
//	return func(db *DB) {
//
//	}
//}

//func MustNewDB(opts ...DBOption) *DB {
//	db, err := NewDB(opts...)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
