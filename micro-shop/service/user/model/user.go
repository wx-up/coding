package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;" `
}

type CommonTimestampsField struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDeleteTimestampField struct {
	// DeletedAt time.Time `gorm:"column:deleted_at;index;"`
	DeletedAt gorm.DeletedAt
}

type User struct {
	BaseModel
	CommonTimestampsField
	SoftDeleteTimestampField

	Phone    string     `gorm:"column:phone;type:varchar(11);not null;unique;index:idx_phone;comment:手机号;default:'';"`
	Password string     `gorm:"column:password;type:varchar(32);not null;comment:密码;default:'';"`
	Nickname string     `gorm:"column:nickname;type:varchar(20);not null;comment:昵称;default:'';"`
	Birthday *time.Time `gorm:"column:birthday;type:datetime;comment:生日;default:NULL;"`

	// 使用字符串的话，比较直白，但是检索效率相比数字会比较差，内存占用也多
	Gender int8 `gorm:"column:gender;type:tinyint(1);not null;comment:性别1男2女;default:1;"`

	Role int8 `gorm:"column:role;type:tinyint(1);not null;comment:角色1普通用户2管理员;default:1;"`
}
