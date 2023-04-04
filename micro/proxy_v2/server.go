package proxy_v2

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
)

type Server struct {
	// 维护服务信息
	services map[string]Service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service, 16),
	}
}

func (s *Server) Register(srv Service) {
	s.services[srv.Name()] = srv
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
		// 先读长度
		lengthBytes := make([]byte, numOfLengthBytes)
		_, err := conn.Read(lengthBytes)
		if err != nil {
			return err
		}
		length := binary.BigEndian.Uint64(lengthBytes)

		// 读取全部的数据
		reqMsg := make([]byte, length)
		_, err = conn.Read(reqMsg)
		if err != nil {
			return err
		}

		data, err := s.handleMsg(reqMsg)
		if err != nil {
			// err 可能是业务错误，需要考虑返回给客户端
		}

		res := make([]byte, numOfLengthBytes, numOfLengthBytes+len(data))
		binary.BigEndian.PutUint64(res, uint64(len(data)))
		res = append(res, data...)
		_, err = conn.Write(res)
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

	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("你要调用的服务不存在")
	}

	// 反射发起调用
	refVal := reflect.ValueOf(service)
	// 参数：一个 context 一个 req
	method := refVal.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	_ = json.Unmarshal(req.Arg, inReq.Interface())
	in[1] = inReq
	result := method.Call(in)
	// 返回值：一个是 result 一个是 error
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}
	res, err = json.Marshal(result[0].Interface())
	return
}
