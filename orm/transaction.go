package orm

import (
	"context"
	"database/sql"

	"github.com/wx-up/coding/orm/internal/model"
	"github.com/wx-up/coding/orm/internal/valuer"
)

// Tx 对事务具柄进行封装（ 有了封装后续才可以为所欲为！）
type Tx struct {
	// 这里不使用组合是为了不想让用户绕开 Tx 这个接口
	// 如果使用了组合，那么 Commit、Rollback 这种方法就不需要定义了
	tx *sql.Tx

	core core
}

func (tx *Tx) getCore() core {
	return tx.core
}

func (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) RollbackIfNotCommit() error {
	// 当事务已经结束（ 提交或者回滚 ），再调用 rollback 会返回 sql.ErrTxDone 错误
	err := tx.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

type Session interface {
	// 直接返回一个包装结构体 core，这样子的好处可以防止接口 Session 膨胀
	// 如果将 Registry 和 Factory 平铺出来就会定义类似如下的接口
	//  getModel() *model.Model
	//  getValCreator() valuer.Factory
	// 这种方式，如果需要返回 DB 的话，Session 接口又要新增一个方法：
	//  getDB() *DB
	// 而 core 的方式只是在 core 结构体中新增一个字段，Session 接口不动，较稳定
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type core struct {
	r          *model.Registry
	valCreator valuer.Factory
	dialect    Dialect
	ms         []Middleware
}

func (c core) Middlewares() []Middleware {
	return c.ms
}

func (c core) Dialect() Dialect {
	return c.dialect
}

func (c core) R() *model.Registry {
	return c.r
}

func (c core) ValCreator() valuer.Factory {
	return c.valCreator
}
