package main

import (
	"context"
	"fmt"
	"net"

	"github.com/wx-up/coding/micro/hello-world/grpc/proto"
	"google.golang.org/grpc"
)

type HelloService struct{}

func (h *HelloService) Hello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{Replay: "Hello," + request.Name}, nil
}

func main() {
	srv := grpc.NewServer()
	proto.RegisterHelloServiceServer(srv, &HelloService{})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		fmt.Println(err)
		return
	}
}
