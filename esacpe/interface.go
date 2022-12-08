package esacpe

import "fmt"

/*
	go build -gcflags="-m" interface.go
	以接口作为返回值或者 interface{} 作为参数或者返回值，都会导致逃逸
*/

func Print() {
	a := 10
	fmt.Println(a)
}

type Person struct{}

func (u Person) Say() {
}

func NewSayer() Sayer {
	return Person{}
}

type Sayer interface {
	Say()
}
