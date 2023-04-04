package response

import (
	"time"

	"github.com/wx-up/coding/micro-shop/service/user/proto"
)

type JsonTime struct {
	time.Time
}

func (t *JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format("2006-01-02 15:04:05")), nil
}

func (t *JsonTime) UnmarshalJSON([]byte) error {
	return nil
}

type UserItem struct {
	Id       uint64   `json:"id"`
	Phone    string   `json:"phone"`
	Nickname string   `json:"nickname"`
	Birthday JsonTime `json:"birthday"`
	Gender   int8     `json:"gender"`
}

func ToUserItem(obj *proto.UserItem) UserItem {
	return UserItem{
		Id:       obj.Id,
		Phone:    obj.Phone,
		Nickname: obj.Nickname,
		Birthday: JsonTime{Time: time.Unix(int64(obj.Birthday), 0)},
		Gender:   int8(obj.Gender),
	}
}

func ToUserItems(objs []*proto.UserItem) []UserItem {
	var items []UserItem
	for _, obj := range objs {
		items = append(items, ToUserItem(obj))
	}
	return items
}
