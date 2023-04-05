package message

import (
	"bytes"
	"encoding/binary"
)

type Request struct {
	HeaderLength uint32 // 协议头长度
	BodyLength   uint32 // 协议体长度
	RequestID    uint32 // 消息ID
	Version      uint8  // 版本
	Compress     uint8  // 压缩算法
	Serializer   uint8  // 序列化协议

	ServiceName string
	MethodName  string

	// value 定义成 any 的话编解码是比较麻烦的
	Meta map[string]string

	// 请求参数
	Data []byte
}

func (req *Request) CalculateHeaderLength() {
	// 长度 1 表示分隔符
	headerLength := 15 + uint32(len(req.ServiceName)) + 1 + uint32(len(req.MethodName)) + 1

	for key, value := range req.Meta {
		headerLength += uint32(len(key))
		headerLength += 1 // key value 之间的分隔符
		headerLength += uint32(len(value))
		headerLength += 1 // key value 和下一个 key value 之间的分隔符
	}

	req.HeaderLength = headerLength
}

func (req *Request) CalculateBodyLength() {
	req.BodyLength = uint32(len(req.Data))
}

func EncodeReq(req *Request) []byte {
	bs := make([]byte, req.BodyLength+req.HeaderLength)

	cur := bs
	binary.BigEndian.PutUint32(cur[:4], req.HeaderLength) // 写入 头部长度

	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], req.BodyLength) // 写入 内容长度

	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], req.RequestID) // 写入 消息ID

	cur = cur[4:]
	cur[0] = req.Version // 版本本身就是一个字节，所以直接赋值即可
	cur = cur[1:]
	cur[0] = req.Compress // 同理
	cur = cur[1:]
	cur[0] = req.Serializer // 同理

	// 追加 server name
	cur = cur[1:]
	serviceNameLength := len(req.ServiceName)
	copy(cur[:serviceNameLength], req.ServiceName)

	// 因为 server name 和 method name 是变长的，我们不知道边界，这里引入 分隔符，比如： user-server\nGetById
	// 后续在解码的时候需要使用这个分隔符
	cur[serviceNameLength] = '\n'

	// 追加 method name
	cur = cur[serviceNameLength+1:]
	methodNameLength := len(req.MethodName)
	copy(cur[:methodNameLength], req.MethodName)

	// 同理，method name 和 meta 之间也需要使用分隔符区分
	cur[methodNameLength] = '\n'
	cur = cur[methodNameLength+1:]

	for key, value := range req.Meta {
		keyLength := len(key)
		copy(cur, key)
		// key 和 value 之间也需要分隔符区分
		cur[keyLength] = '\r'
		cur = cur[keyLength+1:]

		valueLength := len(value)
		copy(cur, value)

		// key value 和下一个 key value 之间的分隔符
		cur[valueLength] = '\n'

		cur = cur[valueLength+1:]
	}

	copy(cur, req.Data)

	return bs
}

func DecodeReq(data []byte) *Request {
	req := &Request{}
	req.HeaderLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.RequestID = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compress = data[13]
	req.Serializer = data[14]

	header := data[15:req.HeaderLength]

	// 分隔符
	index := bytes.IndexByte(header, '\n')

	req.ServiceName = string(header[:index])

	// index 是分隔符本身，所以需要 +1 跳掉分隔符
	header = header[index+1:]
	index = bytes.IndexByte(header, '\n')
	req.MethodName = string(header[:index])

	// meta 以及 meta 之后的部分
	header = header[index+1:]
	index = bytes.IndexByte(header, '\n')
	if index != -1 { // 存在 meta
		// 设置一个预估容量
		meta := make(map[string]string, 4)
		for index != -1 {
			pair := header[:index]

			pairIndex := bytes.IndexByte(pair, '\r')
			key := string(pair[:pairIndex])
			value := string(pair[pairIndex+1:])
			meta[key] = value

			header = header[index+1:]
			index = bytes.IndexByte(header, '\n')
		}
		req.Meta = meta
	}

	if req.BodyLength > 0 {
		req.Data = data[req.HeaderLength:]
	}

	return req
}
