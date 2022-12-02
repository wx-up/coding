package option

type GenericOption[T any] func(*T)

func Apply[T any](t *T, opts ...GenericOption[T]) {
	for _, opt := range opts {
		opt(t)
	}
}

type Dog struct {
	name  string
	hobby string
}

func WithHobby(hobby string) GenericOption[Dog] {
	return func(d *Dog) {
		d.hobby = hobby
	}
}

func NewDog(opts ...GenericOption[Dog]) *Dog {
	d := &Dog{}
	Apply[Dog](d, opts...)
	return d
}
