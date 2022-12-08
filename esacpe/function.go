package esacpe

/*
	go build -gcflags="-m" function.go
	闭包会导致逃逸
*/

func Counter() func() int64 {
	num := int64(0)
	// 闭包
	return func() int64 {
		return num
	}
}
