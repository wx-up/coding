package reflect

import (
	"errors"
	"reflect"
)

// iterateVal 讲结构体的字段和值打印出来
// TDD 思路编写
func iterateVal(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("val 不能为 nil")
	}

	refT := reflect.TypeOf(val)
	refV := reflect.ValueOf(val)

	// for 处理 val 为多级指针的情况：***User{}
	for refT.Kind() == reflect.Pointer {
		refT = refT.Elem()
		refV = refV.Elem()
	}

	if refT.Kind() != reflect.Struct {
		return nil, errors.New("非法输入")
	}

	numsField := refT.NumField()
	res := make(map[string]any)
	for i := 0; i < numsField; i++ {
		// 忽略结构体中字段类型为结构体或者指针的字段
		if refV.Field(i).Kind() == reflect.Pointer || refV.Field(i).Kind() == reflect.Struct {
			continue
		}

		// Go 中的反射只能拿到私有字段的类型信息，拿不到私有字段的值信息 这和 java 有区别
		if refT.Field(i).IsExported() {
			res[refT.Field(i).Name] = refV.Field(i).Interface()
		} else {
			res[refT.Field(i).Name] = reflect.Zero(refT.Field(i).Type).Interface()
		}
	}

	return res, nil
}
