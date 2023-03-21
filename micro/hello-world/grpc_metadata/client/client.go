package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/wx-up/coding/micro/hello-world/grpc_metadata/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct{}

type ClientCredential struct{}

func (c *ClientCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": "123",
	}, nil
}

func (c *ClientCredential) RequireTransportSecurity() bool {
	return false
}

func main() {
	client, err := grpc.Dial("localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&ClientCredential{}),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if client != nil {
		defer func() {
			_ = client.Close()
		}()
	}

	// 填充 metadata
	md := metadata.New(map[string]string{
		"name": "Bob",
		"age":  "18",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := proto.NewUserClient(client).Say(ctx, &proto.SayReq{
		Name: "哈哈哈",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Ok)
}
