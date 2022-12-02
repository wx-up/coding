package pool

import (
	"fmt"
	"sync"
	"testing"
)

func Test(t *testing.T) {
	type User struct {
		Name string
	}

	pool := NewPool[*User](func() *User {
		return &User{}
	})

	obj := pool.Get()
	obj.Name = "wx"
	fmt.Println(obj)
}

type Number interface {
	int | int64
}

func Add[T Number](a, b T) T {
	return a + b
}

func TestAdd(t *testing.T) {
	fmt.Println(Add(1, 2))
}

// BenchmarkPool_Get go test -run=None -bench=.
func BenchmarkPool_Get(b *testing.B) {
	newPool := NewPool[string](func() string {
		return ""
	})

	pool := sync.Pool{
		New: func() any {
			return ""
		},
	}

	b.Run("pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newPool.Get()
		}
	})

	b.Run("sync.pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pool.Get()
		}
	})
}
