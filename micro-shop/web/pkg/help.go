package pkg

import (
	"net"
)

func GetFreePort() (int, error) {
	// 动态获取可用端口
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.Listen("tcp", addr.String()) // addr.String()  ==> 127.0.0.1:0
	if err != nil {
		return 0, err
	}

	return l.Addr().(*net.TCPAddr).Port, nil
}
