package handler

import (
	"github.com/wx-up/coding/micro-shop/service/user/model"
	"github.com/wx-up/coding/micro-shop/service/user/proto"
)

func FormatUsersToUserItems(objs []model.User) []*proto.UserItem {
	if len(objs) <= 0 {
		return nil
	}
	res := make([]*proto.UserItem, 0, len(objs))
	for _, obj := range objs {
		res = append(res, FormatUserToUserItem(obj))
	}
	return res
}

func FormatUserToUserItem(obj model.User) *proto.UserItem {
	res := &proto.UserItem{
		Id:       uint64(obj.ID),
		Phone:    obj.Phone,
		Password: obj.Password,
		Nickname: obj.Nickname,
		Gender:   int32(obj.Gender),
		Role:     int32(obj.Role),
	}
	if obj.Birthday != nil {
		res.Birthday = uint64(obj.Birthday.Unix())
	}
	return res
}
