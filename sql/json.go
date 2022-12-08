package sql

import (
	"database/sql/driver"
	"encoding/json"
)

type JsonColumn struct {
	Val any
	// 为 true 时，Val 才有意义
	Valid bool
}

func (j *JsonColumn) Scan(src any) error {
	if src == nil {
		return nil
	}
	bs := src.([]byte)
	if len(bs) <= 0 {
		return nil
	}
	err := json.Unmarshal(bs, &j.Val)
	j.Valid = err == nil
	return err
}

func (j *JsonColumn) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}
