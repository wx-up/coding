package builder

// UserBuilder 建造者模式，一般是基于链式调用
type UserBuilder struct {
	name string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (ub *UserBuilder) BuildName(name string) *UserBuilder {
	ub.name = name
	return ub
}

type User struct {
	name string
}

func (ub *UserBuilder) Build() *User {
	u := &User{
		name: ub.name,
	}
	return u
}
