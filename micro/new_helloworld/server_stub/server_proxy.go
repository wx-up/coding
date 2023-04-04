package server_proxy

import (
	"net/rpc"

	"github.com/wx-up/coding/micro/new_helloworld/hanlder"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

// 如果做到解耦 - 我们关系的是函数 鸭子类型
func RegisterHelloService(srv HelloServicer) error {
	return rpc.RegisterName(hanlder.HelloServiceName, srv)
}
