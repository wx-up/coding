package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"github.com/wx-up/coding/micro/hello-world/grpc_metadata/proto"
)

type UserService struct{}

func (u *UserService) Say(ctx context.Context, req *proto.SayReq) (*proto.SayResp, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range meta {
			fmt.Println(fmt.Sprintf("%s = %s", k, v))
		}
	}
	fmt.Println((*req).Name)
	return &proto.SayResp{
		Ok: true,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 一元拦截器
	interceptorOpt := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// meta, ok := metadata.FromIncomingContext(ctx)
		return handler(ctx, req)
	})
	srv := grpc.NewServer(interceptorOpt)
	proto.RegisterUserServer(srv, &UserService{})

	err = srv.Serve(listener)
	if err != nil {
		fmt.Println(err)
		return
	}
}
