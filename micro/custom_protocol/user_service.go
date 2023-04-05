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
