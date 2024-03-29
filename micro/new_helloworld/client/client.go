package main

import (
	"fmt"

	"github.com/wx-up/coding/micro/new_helloworld/client_stub"
)

/*
 引入了客户端（ client_stub ）和服务端存根之后，我们就可以更加聚焦业务逻辑的编写
*/

func main() {
	// 1. 建立连接
	client := client_proxy.NewHelloServiceClient("tcp", "localhost:1234")
	// 1. 只想写业务逻辑 不想关注每个函数的名称
	// 客户端部分
	var reply string // string有默认值
	err := client.Hello("bobby", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)

	// 1. 这些概念在grpc中都有对应
	// 2. 发自灵魂的拷问： server_proxy 和 client_proxy能否自动生成啊 为多种语言生成
	// 3. 都能满足 这个就是protobuf + grpc
}
