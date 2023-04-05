package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncodeDecode 如果两个方法是对称的过程，那么可以使用下面的方式一起测试。类似的场景：加解密等等
// 其次这种对称的函数是非常适合模糊测试的
func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		Name string
		req  *Request
	}{
		{
			Name: "normal",
			req: &Request{
				RequestID:   123,
				Version:     1,
				Compress:    1,
				Serializer:  1,
				ServiceName: "user-server",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123456",
					"a/b":      "a",
				},
				Data: []byte("hello world"),
			},
		},
		{
			// data 有换行符是很正常的，比如文章内容
			// meta 可以禁用协议的分隔符，比如 \n 和 \r
			Name: "data with \n",
			req: &Request{
				RequestID:   123,
				Version:     1,
				Compress:    1,
				Serializer:  1,
				ServiceName: "user-server",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123456",
					"a/b":      "a",
				},
				Data: []byte("hello \n world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.req.CalculateBodyLength()
			tc.req.CalculateHeaderLength()
			data := EncodeReq(tc.req)
			res := DecodeReq(data)
			assert.Equal(t, tc.req, res)
		})
	}
}
