package handler

import (
	"context"
	"crypto/sha512"
	"fmt"

	"github.com/wx-up/coding/micro-shop/service/user/pkg/password"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/wx-up/coding/micro-shop/service/user/model"

	"github.com/wx-up/coding/micro-shop/service/user/global"

	"github.com/wx-up/coding/micro-shop/service/user/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct{}

// CheckPassword 检测密码是否相等（ 登录的时候需要 ）
func (u *UserService) CheckPassword(ctx context.Context, req *proto.CheckReq) (*proto.CheckResp, error) {
	res := &proto.CheckResp{IsEqual: false}
	if generatePassword(req.Password) == req.EncodePassword {
		res.IsEqual = true
	}
	return res, nil
}

func (u *UserService) Users(ctx context.Context, req *proto.ListReq) (*proto.ListResp, error) {
	if err := ValidateListReq(req); err != nil {
		return nil, err
	}

	var count int64
	if err := global.DB().Model(model.User{}).Count(&count).Error; err != nil {
		return nil, err
	}

	var objs []model.User
	if err := global.DB().Limit(int(req.PageSize)).
		Offset(int((req.PageIndex - 1) * req.PageSize)).
		Find(&objs).Error; err != nil {
		return nil, err
	}
	return &proto.ListResp{
		Total: uint32(count),
		Items: FormatUsersToUserItems(objs),
	}, nil
}

func (u *UserService) GetUserByPhone(ctx context.Context, req *proto.PhoneReq) (*proto.UserItem, error) {
	var obj model.User
	if err := global.DB().Model(model.User{}).Where("phone = ?", req.Phone).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		return nil, err
	}
	return FormatUserToUserItem(obj), nil
}

func (u *UserService) GetUserById(ctx context.Context, req *proto.IdReq) (*proto.UserItem, error) {
	var obj model.User
	if err := global.DB().Model(model.User{}).Where("id = ?", req.Id).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		return nil, err
	}
	return FormatUserToUserItem(obj), nil
}

func (u *UserService) CreateUser(ctx context.Context, req *proto.CreateReq) (*proto.UserItem, error) {
	if err := ValidateCreateReq(req); err != nil {
		return nil, err
	}
	// 查询
	obj, err := findUser("phone", req.Phone)
	if err != nil {
		return nil, err
	}
	if obj.ID > 0 {
		return nil, status.Error(codes.AlreadyExists, "用户已经存在")
	}

	obj.Phone = req.Phone
	obj.Role = int8(req.Role)

	// 加密密码
	obj.Password = generatePassword(req.Password)

	// 创建用户
	if err = global.DB().Create(&obj).Error; err != nil {
		return nil, err
	}
	return FormatUserToUserItem(obj), nil
}

func (u *UserService) UpdateUser(ctx context.Context, req *proto.UpdateReq) (*emptypb.Empty, error) {
	// 查询
	obj, err := findUser("id", req.Id)
	if err != nil {
		return nil, err
	}
	if obj.ID <= 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	// 更新
	columns := generateUpdateColumns(req)
	if len(columns) <= 0 {
		return nil, nil
	}
	return nil, global.DB().Model(model.User{}).Where("id", obj.ID).Updates(columns).Error
}

func findUser(column string, value interface{}) (model.User, error) {
	var obj model.User
	err := global.DB().Model(model.User{}).Where(fmt.Sprintf("%s = ?", column), value).First(&obj).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return model.User{}, err
	}

	return obj, nil
}

func generateUpdateColumns(req *proto.UpdateReq) map[string]interface{} {
	columns := make(map[string]interface{})
	if req.Password != "" {
		columns["password"] = generatePassword(req.Password)
	}
	if req.Birthday > 0 {
	}
	if req.Nickname != "" {
		columns["nickname"] = req.Nickname
	}
	if req.Role > 0 {
		if !(req.Role == 1 || req.Role == 2) {
			req.Role = 1
		}
		columns["role"] = req.Role
	}
	if req.Gender > 0 {
		if !(req.Gender == 1 || req.Gender == 2) {
			req.Gender = 1
		}
		columns["gender"] = req.Gender
	}
	if req.Phone != "" {
		columns["phone"] = req.Phone
	}
	return columns
}

func generatePassword(pwd string) string {
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encodePassword := password.Encode(pwd, &options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePassword)
}
