package main

import (
	"context"
	"fmt"

	"github.com/wx-up/coding/micro/hello-world/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	if client != nil {
		defer func() {
			_ = client.Close()
		}()
	}

	res, err := proto.NewHelloServiceClient(client).Hello(context.Background(), &proto.HelloRequest{
		Name: "哈哈",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Replay)
}
