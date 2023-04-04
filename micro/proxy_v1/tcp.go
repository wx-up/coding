package proxy_v1

import "encoding/binary"

// lengthBytes 用多少个字节来表达长度
const lengthBytes = 8

func EncodeMsg(data []byte) []byte {
	res := make([]byte, len(data)+lengthBytes)

	// 大顶端编码，把长度编码成二进制，然后放到了 resp 的前八个字节
	binary.BigEndian.PutUint64(res, uint64(len(data)))
	copy(res[lengthBytes:], data)
	return res
}
