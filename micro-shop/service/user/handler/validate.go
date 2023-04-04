package handler

import "github.com/wx-up/coding/micro-shop/service/user/proto"

func ValidateListReq(req *proto.ListReq) error {
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	return nil
}

func ValidateCreateReq(req *proto.CreateReq) error {
	if !(req.Role == 1 || req.Role == 2) {
		req.Role = 1
	}
	return nil
}
