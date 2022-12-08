package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonColumn[T any] struct {
	Valid bool // true 时 Val 可用
	Val   T
}

// Value Go类型转成数据库类型的时候会调用Value方法
func (j *JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

// Scan 数据库类型转成Go类型会调用Scan方法
func (j *JsonColumn[T]) Scan(src any) error {
	if src == nil {
		return errors.New("src 不能为 nil")
	}
	var bs []byte
	switch v := src.(type) {
	case string:
		bs = []byte(v)
	case []byte:
		bs = v
	}
	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
