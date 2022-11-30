package builder

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
