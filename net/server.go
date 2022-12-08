package net

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func() {
			err := s.handleConn(conn)
			if err != nil {
				_ = conn.Close()
				fmt.Printf("con error：%v\n", err)
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		// 读数据长度
		bs := make([]byte, lenContent)
		_, err := conn.Read(bs)
		if err != nil {
			return err
		}

		reqBs := make([]byte, binary.BigEndian.Uint64(bs))
		_, err = conn.Read(reqBs)
		if err != nil {
			return err
		}
		res := string(reqBs) + ", from response"
		// 总长度
		bs = make([]byte, lenContent, lenContent+len(res))
		// 写入消息长度
		binary.BigEndian.PutUint64(bs, uint64(len(res)))
		bs = append(bs, res...)
		_, err = conn.Write(bs)
		if err != nil {
			return err
		}
	}
}
