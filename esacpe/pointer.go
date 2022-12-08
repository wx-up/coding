package esacpe

/*
	go build -gcflags="-m" pointer.go
	函数返回指针，则会发生逃逸
*/

type User struct{}

func NewUser() *User {
	return &User{}
}
