package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
使用 json 编码协议进行通信
*/
func main() {
	// 直接采用 net 包进行拨号，得到的是一个 conn
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	var reply string

	// 将 conn 包装成一个 json client
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// 数据以如下的 json 格式发送给服务端：{"method":"", "params":["world"], id:0}
	err = client.Call("HelloService.Hello", "world", &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply)
}
