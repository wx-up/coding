package list

import (
	"fmt"
)

type ArrayList[T any] struct {
	values []T
}

func NewArrayList[T any](cap int) *ArrayList[T] {
	return &ArrayList[T]{
		values: make([]T, 0, cap),
	}
}

func (a *ArrayList[T]) Len() int {
	return len(a.values)
}

func (a *ArrayList[T]) Get(index int) (t T, err error) {
	l := a.Len()
	if index <= 0 || index >= l {
		return t, fmt.Errorf("超出范围，下标：%d，长度：%d", index, l)
	}
	return a.values[index], nil
}

func (a *ArrayList[T]) Append(v ...T) {
	a.values = append(a.values, v...)
}

func (a *ArrayList[T]) Add(index int, v T) error {
	l := a.Len()
	if index < 0 || index >= l {
		return fmt.Errorf("超出范围，下标：%d，长度：%d", index, l)
	}
	return nil
}
