package api

import (
	"errors"

	"github.com/wx-up/coding/micro-shop/web/config"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/coding/micro-shop/service/user/proto"
	"github.com/wx-up/coding/micro-shop/web/pkg/client"
	"github.com/wx-up/coding/micro-shop/web/request"
	"github.com/wx-up/coding/micro-shop/web/response"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) List(ctx *gin.Context) {
	var in request.UserListReq
	if err := request.Validate(ctx, &in); err != nil {
		response.Fail(ctx, err)
		return
	}

	resp, err := client.Get(config.Config().UserServer.Name).MustToUserClient().Users(ctx.Request.Context(), &proto.ListReq{
		PageIndex: uint32(in.PageIndex),
		PageSize:  uint32(in.PageSize),
	})
	if err != nil {
		response.Fail(ctx, err)
		return
	}

	// 处理 rpc 返回的结果
	response.Success(ctx, response.ToUserItems(resp.Items))
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var in request.UserLoginReq
	if err := request.Validate(ctx, &in); err != nil {
		response.Fail(ctx, err)
		return
	}
	resp, err := client.Get(config.Config().UserServer.Name).MustToUserClient().GetUserByPhone(ctx.Request.Context(), &proto.PhoneReq{
		Phone: in.Account,
	})
	if err != nil || resp == nil {
		response.Fail(ctx, errors.New("用户名或者密码错误"))
		return
	}

	// 验证密码是否正确
	checkResp, err := client.Get(config.Config().UserServer.Name).MustToUserClient().CheckPassword(ctx.Request.Context(), &proto.CheckReq{
		Password:       in.Password,
		EncodePassword: resp.Password,
	})
	if err != nil || !checkResp.IsEqual {
		response.Fail(ctx, errors.New("用户名或者密码错误"))
		return
	}

	// 用户名密码正确
}
