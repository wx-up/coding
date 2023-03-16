package main

import (
	"context"
	"fmt"

	"github.com/wx-up/coding/micro/hello-world/grpc_server_stream/proto"

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

	res, err := proto.NewHelloClient(client).GetStream(context.Background(), &proto.HelloReq{
		Data: "test",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		resp, err := res.Recv()
		// 服务端发送结束之后，客户端会收到一个 EOF 的错误
		// 客户端发送结束之后，服务端会收到一个 context canceled 错误
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(resp.Data)
	}
}
