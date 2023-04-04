package client

import (
	"fmt"
	"sync"

	"github.com/wx-up/coding/micro-shop/service/user/proto"

	"google.golang.org/grpc"
)

type Client struct {
	Conn *grpc.ClientConn

	Err error
}

var clients = make(map[string]*Client)

var mutex sync.RWMutex

func Register(name string, conn *grpc.ClientConn) {
	mutex.Lock()
	defer mutex.Unlock()
	clients[name] = &Client{
		Conn: conn,
	}
}

func Get(name string) *Client {
	mutex.RLock()
	defer mutex.RUnlock()
	conn, ok := clients[name]
	if !ok {
		return &Client{
			Err: fmt.Errorf("%s 服务不存在，请先注册", name),
		}
	}
	return conn
}

func Close() {
	for key, client := range clients {
		if client.Conn == nil {
			continue
		}
		if err := client.Conn.Close(); err != nil {
			fmt.Printf("%s 服务关闭失败：%v", key, err)
		}
	}
}

func (c *Client) MustToUserClient() proto.UserClient {
	if c.Err != nil {
		panic(c.Err)
	}
	return proto.NewUserClient(c.Conn)
}
