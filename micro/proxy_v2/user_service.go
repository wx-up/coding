package proxy_v2

import "context"

// 客户端
type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (s UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Msg string
}

// UserServiceServer 服务端
type UserServiceServer struct{}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Msg: "我收到了",
	}, nil
}
