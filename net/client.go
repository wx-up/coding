package net

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Send(msg string) (string, error) {
	conn, err := net.DialTimeout("tcp", c.addr, 3*time.Second)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	// 发送数据
	bs := make([]byte, lenContent, lenContent+len(msg))
	binary.BigEndian.PutUint64(bs, uint64(lenContent))
	bs = append(bs, msg...)
	_, err = conn.Write(bs)
	if err != nil {
		return "", err
	}

	// 读取响应：读取内容长度
	lenBs := make([]byte, lenContent)
	_, err = conn.Read(lenBs)
	if err != nil {
		return "", err
	}
	contentLen := binary.BigEndian.Uint64(lenBs)

	// 读取响应：读取内容
	contentBs := make([]byte, contentLen)
	_, err = conn.Read(contentBs)
	if err != nil {
		return "", err
	}
	return string(contentBs), nil
}
