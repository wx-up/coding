package message

import (
	"encoding/binary"
)

type Response struct {
	HeaderLength uint32 // 协议头长度
	BodyLength   uint32 // 协议体长度
	RequestID    uint32 // 消息ID
	Version      uint8  // 版本
	Compress     uint8  // 压缩算法
	Serializer   uint8  // 序列化协议

	Error []byte // 错误

	Data []byte
}

func EncodeResp(resp *Response) []byte {
	data := make([]byte, resp.HeaderLength+resp.BodyLength)

	cur := data
	binary.BigEndian.PutUint32(cur, resp.HeaderLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur, resp.BodyLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur, resp.RequestID)
	cur = cur[4:]
	cur[0] = resp.Version
	cur = cur[1:]
	cur[0] = resp.Compress
	cur = cur[1:]
	cur[0] = resp.Serializer
	cur = cur[1:]

	copy(cur, resp.Error)

	cur = cur[len(resp.Error):]
	copy(cur, resp.Data)
	return data
}

func DecodeResp(data []byte) *Response {
	resp := &Response{}
	resp.HeaderLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.RequestID = binary.BigEndian.Uint32(data[8:12])
	resp.Version = data[12]
	resp.Compress = data[13]
	resp.Serializer = data[14]

	err := data[15:resp.HeaderLength]
	if len(err) > 0 {
		resp.Error = err
	}

	if resp.BodyLength > 0 {
		resp.Data = data[resp.HeaderLength:]
	}

	return resp
}

func (r *Response) calculateHeaderLength() {
	r.HeaderLength = 4 + 4 + 4 + 1 + 1 + 1 + uint32(len(r.Error))
}

func (r *Response) calculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}
