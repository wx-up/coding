package option

type Option func(*User)

func WithName(name string) Option {
	return func(user *User) {
		user.name = name
	}
}

func WithAge(age int64) Option {
	return func(user *User) {
		user.age = age
	}
}

type User struct {
	id   int64
	name string
	age  int64
}

func NewUser(id int64, opts ...Option) *User {
	u := &User{
		id:   id,
		name: "bob",
		age:  12,
	}

	for _, opt := range opts {
		opt(u)
	}
	return u
}
