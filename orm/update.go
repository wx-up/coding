package orm

import "context"

type Updater[T any] struct {
}

func (u *Updater[T]) Build() (*Query, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	return Result{}
}
