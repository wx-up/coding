package proxy_v2

import (
	"encoding/binary"
	"net"
)

const numOfLengthBytes = 8

func Read(conn net.Conn) ([]byte, error) {
	// 两段读：先读取内容的长度，再根据长度读取内容
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	contentLen := binary.BigEndian.Uint64(lenBs)

	// 读取响应：读取内容
	contentBs := make([]byte, contentLen)
	_, err = conn.Read(contentBs)
	if err != nil {
		return nil, err
	}
	return contentBs, nil
}

func EncodeData(data []byte) []byte {
	res := make([]byte, numOfLengthBytes, numOfLengthBytes+len(data))
	binary.BigEndian.PutUint64(res, uint64(len(data)))
	res = append(res, data...)
	return res
}
