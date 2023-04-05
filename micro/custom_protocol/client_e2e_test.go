package proxy_v2

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockUserServiceServer 服务端
type MockUserServiceServer struct {
	Err error
	Msg string
}

func (u *MockUserServiceServer) Name() string {
	return "user-service"
}

func (u *MockUserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	if u.Err != nil {
		return nil, u.Err
	}
	return &GetByIdResp{
		Msg: u.Msg,
	}, nil
}

func Test(t *testing.T) {
	srv := NewServer()

	// 这里是服务端的 Service 并且需要实现 GetById 方法
	usrSrv := &MockUserServiceServer{}
	srv.Register(usrSrv)
	go func() {
		err := srv.Start(":8088")
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(time.Second * 2)

	// 这里是客户端的 Service 不需要实现 GetById 通过反射篡改成 RPC 调用
	usrClient := &UserService{}
	err := InitClientProxy(":8088", usrClient)
	assert.Equal(t, nil, err)
	if err != nil {
		return
	}

	testCases := []struct {
		name     string
		mock     func()
		wantErr  error
		wantResp any
	}{
		{
			name: "normal",
			mock: func() {
				usrSrv.Err = nil
				usrSrv.Msg = "hello world"
			},
			wantErr:  nil,
			wantResp: &GetByIdResp{Msg: "hello world"},
		},
		{
			name: "error",
			mock: func() {
				usrSrv.Err = errors.New("net fail")
			},
			wantErr: errors.New("net fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, err := usrClient.GetById(context.Background(), &GetByIdReq{
				Id: 123,
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
