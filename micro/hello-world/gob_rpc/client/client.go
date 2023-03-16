package main

import (
	"fmt"
	"net/rpc"
)

/*
 默认采用 gob 协议通信
*/

func main() {
	// 直接使用 rpc 包拨号，得到的是一个 client
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	var reply string
	err = client.Call("HelloService.Hello", "world", &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply)
}
