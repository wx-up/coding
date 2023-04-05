package proxy_v2

import (
	"encoding/binary"
	"net"
)

const numOfLengthBytes = 8

func Read(conn net.Conn) ([]byte, error) {
	// 根据自定义协议的设计，前 8 个字节表示协议头的长度和协议体的长度
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	headerLength := binary.BigEndian.Uint32(lenBs[:4])  // 协议头长度
	contentLength := binary.BigEndian.Uint32(lenBs[4:]) // 协议体长度
	length := headerLength + contentLength
	data := make([]byte, length)
	// 读取其他内容，前面已经读取了8个字节了
	_, err = conn.Read(data[8:])
	if err != nil {
		return nil, err
	}

	// 填充前面已经读取的 8 个字节
	copy(data[:8], lenBs)

	return data, nil
}

func EncodeData(data []byte) []byte {
	res := make([]byte, numOfLengthBytes, numOfLengthBytes+len(data))
	binary.BigEndian.PutUint64(res, uint64(len(data)))
	res = append(res, data...)
	return res
}
