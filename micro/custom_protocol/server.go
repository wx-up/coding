package proxy_v2

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"

	"github.com/wx-up/coding/micro/custom_protocol/message"
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
		resp, err := s.Invoke(context.Background(), message.DecodeReq(reqMsg))
		if err != nil {
			// 反射调用函数出现问题，将这个错误返回给调用方
			resp.Error = []byte(err.Error())
		}
		resp.CalculateBodyLength()
		resp.CalculateHeaderLength()

		_, err = conn.Write(message.EncodeResp(resp))
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	stub, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("你要调用的服务不存在")
	}

	resp := &message.Response{
		RequestID:  req.RequestID,
		Version:    req.Version,
		Compress:   req.Compress,
		Serializer: req.Serializer,
	}

	// 反射发起调用
	respData, err := stub.invoke(context.Background(), req.MethodName, req.Data)
	if err != nil {
		return resp, err
	}
	resp.Data = respData
	return resp, nil
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
