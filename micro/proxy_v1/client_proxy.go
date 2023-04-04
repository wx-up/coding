package proxy_v1

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/silenceper/pool"
)

type Client struct {
	// 这里不使用 sync.Pool 的原因是它不会帮你 close 连接（ 垃圾回收机制 ）
	// 而 pool.Pool 会帮你 close 连接
	pool pool.Pool
}

// NewClient addr 这个参数可测试性不是很好
// 传入 factory 可测试性会更好点
func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		// 压测的时候，主要是调整 InitialCap 和 MaxCap 这两个值
		// 有可能性能慢是因为一开始应用没有较多的初始化连接，导致前面的请求特别慢，那就可以调整  InitialCap 的值
		// 有可能性能慢的另一个原因就是连接数量不够，那就可以调整 MaxCap 的值
		// 压测的时候，一般用不上这个 MaxIdle 参数 因为在压测的时候，不会有连接是空闲的状态
		InitialCap: 10,  // 初始容量
		MaxCap:     100, // 最多能有多少
		MaxIdle:    20,  // 最大的空闲数量

		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		Close: func(v interface{}) error {
			return v.(net.Conn).Close()
		},

		// 空闲 1 分钟就关闭
		IdleTimeout: time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		pool: p,
	}, nil
}

// Invoke 发送数据需要分为两部分
// 第一部分是长度，内容的长度
// 第二部分才是内容
// 读取方，先读取长度，接着读取指定长度的内容即可。
func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 把数据发送出去
	obj, err := c.pool.Get()
	if err != nil {
		return nil, err
	}

	data = EncodeMsg(data)
	l, err := obj.(net.Conn).Write(data)
	if err != nil {
		return nil, err
	}
	if l != len(data) {
		return nil, errors.New("micro：未写入全部的数据")
	}

	// 处理收到的数据

	return &Response{Data: []byte(`{"name":"123"}`)}, nil
}
