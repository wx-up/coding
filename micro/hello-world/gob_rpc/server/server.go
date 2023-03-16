package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func (s *HelloService) Name() string {
	return "HelloService"
}

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := &HelloService{}
	err = rpc.RegisterName(srv.Name(), srv)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 默认使用的是 Gob 协议
		// 可以修改成其他协议：json 协议
		//  rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		go rpc.ServeConn(conn)
	}
}
