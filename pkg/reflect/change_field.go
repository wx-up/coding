package reflect

import (
	"errors"
	"reflect"
)

// SetField 反射修改字段值
func SetField(entity any, field string, newVal any) error {
	if entity == nil {
		return errors.New("参数错误：entity 为 nil")
	}
	refV := reflect.ValueOf(entity)
	if !(refV.Kind() == reflect.Ptr && refV.Elem().Kind() == reflect.Struct) {
		return errors.New("参数错误：entity 只能是指向结构体的指针")
	}

	refV = refV.Elem()
	refT := refV.Type()

	// Type 的 FieldByName 可以判断出字段是否存在
	// Value 的 FieldByName 则不行
	if _, found := refT.FieldByName(field); !found {
		return errors.New("不存在的字段")
	}

	f := refV.FieldByName(field)

	if !f.CanSet() {
		return errors.New("该字段不可修改")
	}

	f.Set(reflect.ValueOf(newVal))
	return nil
}
