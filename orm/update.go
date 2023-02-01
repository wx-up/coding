package orm

import (
	"context"
)

// Updater 与 upsert 语句中更新部分一样
type Updater[T any] struct {
}

// Where 条件
func (u *Updater[T]) Where() *Updater[T] {
	return u
}

// Set 赋值语句
func (u *Updater[T]) Set() *Updater[T] {
	return u
}

func (u *Updater[T]) Build() (*Query, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	return Result{}
}
