package proxy_v2

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/wx-up/coding/micro/custom_protocol/serialize/json"

	"github.com/wx-up/coding/micro/custom_protocol/message"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_setFuncField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name string

		service Service

		wantErr error

		mock func() Proxy
	}{
		{
			name:    "nil",
			service: nil,
			wantErr: errors.New("rpc：不支持 nil"),
			mock: func() Proxy {
				return NewMockProxy(ctrl)
			},
		},
		{
			name:    "非指向结构体的指针",
			service: UserService{},
			wantErr: errors.New("rpc：只支持指向结构体的指针"),
			mock: func() Proxy {
				return NewMockProxy(ctrl)
			},
		},
		{
			name:    "user service",
			service: &UserService{},
			mock: func() Proxy {
				res := NewMockProxy(ctrl)
				// 期望调用 Invoke 方法时，第一个参数值是任意的不限制，第二个参数值必须是下面指定的
				res.EXPECT().Invoke(gomock.Any(), &message.Request{
					ServiceName: "user-service",
					MethodName:  "GetById",
				}).Return(&message.Response{}, nil)
				return res
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setFuncField(tc.service, tc.mock(), json.New())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: 1})
			require.NoError(t, err)
			fmt.Println(resp)
		})
	}
}
