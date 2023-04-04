package proxy_v2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	srv := NewServer()

	// 这里是服务端的 Service 并且需要实现 GetById 方法
	usrSrv := &UserServiceServer{}
	srv.Register(usrSrv)
	go func() {
		err := srv.Start(":8088")
		if err != nil {
			panic(err)
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

	resp, err := usrClient.GetById(context.Background(), &GetByIdReq{
		Id: 123,
	})
	assert.Equal(t, nil, err)
	if err != nil {
		return
	}
	fmt.Println(resp.Msg)
}
