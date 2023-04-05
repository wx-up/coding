package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseEncodeDecode(t *testing.T) {
	testCases := []struct {
		Name string
		Resp *Response
	}{
		{
			Name: "normal",
			Resp: &Response{
				RequestID:  123,
				Version:    1,
				Compress:   1,
				Serializer: 1,
			},
		},
		{
			Name: "with body",
			Resp: &Response{
				RequestID:  123,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Data:       []byte("123"),
			},
		},
		{
			Name: "with error",
			Resp: &Response{
				RequestID:  123,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Error:      []byte("123"),
			},
		},
		{
			Name: "with error and with body",
			Resp: &Response{
				RequestID:  123,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Error:      []byte("123"),
				Data:       []byte("123"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Resp.calculateBodyLength()
			tc.Resp.calculateHeaderLength()
			data := EncodeResp(tc.Resp)
			res := DecodeResp(data)
			assert.Equal(t, tc.Resp, res)
		})
	}
}
