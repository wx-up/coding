package reflect

import (
	"errors"
	"reflect"
)

func Copy(src any, dst any, ignoreFields []string) error {
	if src == nil {
		return errors.New("src 不能为 nil")
	}
	if dst == nil {
		return errors.New("dst 不能为 nil")
	}
	srcVal := reflect.ValueOf(src)
	srcType := srcVal.Type()
	if !(srcType.Kind() == reflect.Pointer && srcType.Elem().Kind() == reflect.Struct) {
		return errors.New("src 只能是指向结构体的指针")
	}
	srcVal = srcVal.Elem()
	srcType = srcType.Elem()

	dstVal := reflect.ValueOf(dst)
	dstType := dstVal.Type()
	if !(dstType.Kind() == reflect.Pointer && dstType.Elem().Kind() == reflect.Struct) {
		return errors.New("dst 只能是指向结构体的指针")
	}
	dstVal = dstVal.Elem()
	dstType = dstType.Elem()

	srcFieldNum := srcType.NumField()

	// 遍历 src 的所有字段
	for i := 0; i < srcFieldNum; i++ {
		fieldName := srcType.Field(i).Name
		fieldValue := srcVal.Field(i)

		// 在 dst 中查找对应的字段
		_, found := dstType.FieldByName(fieldName)
		if !found {
			continue
		}

		// 找到则设置值
		foundValue := dstVal.FieldByName(fieldName)
		if foundValue.CanSet() {
			foundValue.Set(fieldValue)
		}
	}

	return nil
}
