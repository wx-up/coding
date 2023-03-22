package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 一元拦截器
	interceptorOpt := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if validate, ok := req.(interface{ Validate() error }); ok {
			if err = validate.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}
		// meta, ok := metadata.FromIncomingContext(ctx)
		return handler(ctx, req)
	})
	srv := grpc.NewServer(interceptorOpt)

	err = srv.Serve(listener)
	if err != nil {
		fmt.Println(err)
		return
	}
}
