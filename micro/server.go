package micro

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

type Server struct{}

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
		lengthBytes := make([]byte, lengthBytes)
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

		// 解码
		req := &Request{}
		err = json.Unmarshal(reqMsg, req)
		if err != nil {
			return err
		}

		log.Println(req)
		// 处理请求

	}
}
