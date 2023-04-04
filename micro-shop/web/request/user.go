package request

import "github.com/thedevsaddam/govalidator"

type UserListReq struct {
	PageIndex int64 `json:"page_index" form:"page_index"`
	PageSize  int64 `json:"page_size" form:"page_size"`
}

func (r *UserListReq) Validate() map[string][]string {
	if r.PageIndex <= 0 {
		r.PageIndex = 1
	}
	if r.PageSize <= 0 || r.PageIndex > 100 {
		r.PageSize = 100
	}
	return nil
}

type UserLoginReq struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
}

func (r *UserLoginReq) Validate() map[string][]string {
	rules := govalidator.MapData{
		"account":  []string{"required", "digits:11"},
		"password": []string{"required"},
	}
	messages := govalidator.MapData{
		"account": []string{
			"required:请填写账号",
			"digits:手机号格式不正确",
		},
		"password": []string{
			"required:请填写密码",
		},
	}

	return govalidator.New(govalidator.Options{
		Data:     r,
		Rules:    rules,
		Messages: messages,
	}).ValidateJSON()
}
