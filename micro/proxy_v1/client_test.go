package proxy_v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int64
}

type GetByIdResp struct {
	Name string `json:"name"`
}

// mockProxy
// 我们可以通过 gomock 来创建 mock 实现
// 还有一种方式：就是创建一个 mock struct 实现要 mock 的接口即可
type mockProxy struct {
	req    *Request
	result []byte
	err    error
}

func (m *mockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	m.req = req
	return &Response{Data: m.result}, m.err
}

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name    string
		service *UserService

		// *Request 结构本来是不好断言判断的
		// 我们借助 mockProxy 来对生成的 Request 断言
		proxy *mockProxy

		in       *GetByIdReq
		wantReq  *Request
		wantResp *GetByIdResp

		// InitClientProxy 方法返回 error
		wantInitErr error

		// 调用 GetById 返回 error
		wantErr error
	}{
		{
			name:    "success",
			service: &UserService{},
			proxy: &mockProxy{
				result: []byte(`{"name":"wx"}`),
			},
			wantReq: &Request{
				ServerName: "user-service",
				MethodName: "GetById",
				Arg:        &GetByIdReq{Id: 1},
			},
			wantResp: &GetByIdResp{Name: "wx"},
			in: &GetByIdReq{
				Id: 1,
			},
			wantErr: nil,
		},
		{
			name:    "proxy error",
			service: &UserService{},
			proxy: &mockProxy{
				err: errors.New("请求失败"),
			},
			wantErr: errors.New("请求失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service, tc.proxy)
			assert.Equal(t, tc.wantInitErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.GetById(context.Background(), tc.in)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantReq, tc.proxy.req)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
