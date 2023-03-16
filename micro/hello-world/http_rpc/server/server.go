package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
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
	srv := &HelloService{}
	err := rpc.RegisterName(srv.Name(), srv)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Handle("/httprpc", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			writer,
			request.Body,
		}
		fmt.Println(rpc.ServeRequest(jsonrpc.NewServerCodec(conn)))
	}))
	http.ListenAndServe(":8080", nil)
}
