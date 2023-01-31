package orm

import "context"

// Updater 与 upsert 语句中更新部分一样
type Updater[T any] struct {
}

func (u *Updater[T]) Build() (*Query, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	return Result{}
}
