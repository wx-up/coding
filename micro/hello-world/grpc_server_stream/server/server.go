package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/wx-up/coding/micro/hello-world/grpc_server_stream/proto"
)

type Server struct{}

func (s *Server) GetStream(req *proto.HelloReq, server proto.Hello_GetStreamServer) error {
	for i := 0; i < 10; i++ {
		_ = server.Send(&proto.HelloResp{Data: fmt.Sprintf("%d", time.Now().UnixNano())})
		time.Sleep(time.Second)
	}
	return nil
}

func (s *Server) PutStream(server proto.Hello_PutStreamServer) error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) AllStream(server proto.Hello_AllStreamServer) error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		// 接收
		res, err := server.Recv()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.Data)
	}()
	go func() {
		defer waitGroup.Done()
		// 发送
		err := server.Send(&proto.HelloResp{Data: "发送"})
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	waitGroup.Wait()
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := grpc.NewServer()
	proto.RegisterHelloServer(srv, &Server{})
	err = srv.Serve(listener)
	if err != nil {
		fmt.Println(err)
		return
	}
}
