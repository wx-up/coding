package proxy_v2

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
)

type Server struct {
	// 维护服务信息
	services map[string]*reflectionStub
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]*reflectionStub, 16),
	}
}

func (s *Server) Register(srv Service) {
	s.services[srv.Name()] = &reflectionStub{
		s:     srv,
		value: reflect.ValueOf(srv),
	}
}

func (s *Server) Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("服务已经启动")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// 每一个 conn 由一个 goroutine 去处理
		go func() {
			if err := s.handleConn(conn); err != nil {
				log.Println(err)
				_ = conn.Close()
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqMsg, err := Read(conn)
		if err != nil {
			return err
		}
		data, err := s.handleMsg(reqMsg)
		if err != nil {
			// err 可能是业务错误，需要考虑返回给客户端
		}

		_, err = conn.Write(EncodeData(data))
		if err != nil {
			return err
		}
	}
}

func (s *Server) handleMsg(data []byte) (res []byte, err error) {
	// 解码
	req := &Request{}
	err = json.Unmarshal(data, req)
	if err != nil {
		return
	}

	stub, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("你要调用的服务不存在")
	}

	// 反射发起调用
	return stub.invoke(context.Background(), req.MethodName, req.Arg)
}

// reflectionStub 反射的桩，后续可以考虑用 unsafe 优化
type reflectionStub struct {
	s     Service
	value reflect.Value
}

func (s *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := s.value.MethodByName(methodName)
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	_ = json.Unmarshal(data, inReq.Interface())
	in[1] = inReq
	result := method.Call(in)
	// 返回值：一个是 result 一个是 error
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}
	return json.Marshal(result[0].Interface())
}
