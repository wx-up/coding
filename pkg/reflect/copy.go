package reflect

import (
	"errors"
	"reflect"
)

// CopyBuilder 结构体拷贝，不支持组合等情况
type CopyBuilder struct {
	src          any
	dst          any
	ignoreFields []string
}

func (cb *CopyBuilder) BuildSrc(src any) *CopyBuilder {
	cb.src = src
	return cb
}

func (cb *CopyBuilder) BuildDst(dst any) *CopyBuilder {
	cb.dst = dst
	return cb
}

func (cb *CopyBuilder) BuildIgnoreFields(ignore []string) *CopyBuilder {
	cb.ignoreFields = ignore
	return cb
}

func (cb *CopyBuilder) Builder() error {
	if cb.src == nil {
		return errors.New("src 不能为 nil")
	}
	if cb.dst == nil {
		return errors.New("dst 不能为 nil")
	}
	srcVal := reflect.ValueOf(cb.src)
	srcType := srcVal.Type()
	if !(srcType.Kind() == reflect.Pointer && srcType.Elem().Kind() == reflect.Struct) {
		return errors.New("src 只能是指向结构体的指针")
	}
	srcVal = srcVal.Elem()
	srcType = srcType.Elem()

	dstVal := reflect.ValueOf(cb.dst)
	dstType := dstVal.Type()
	if !(dstType.Kind() == reflect.Pointer && dstType.Elem().Kind() == reflect.Struct) {
		return errors.New("dst 只能是指向结构体的指针")
	}
	dstVal = dstVal.Elem()
	dstType = dstType.Elem()

	srcFieldNum := srcType.NumField()
	for i := 0; i < srcFieldNum; i++ {
		fieldName := srcType.Field(i).Name
		fieldValue := srcVal.Field(i)
		_, found := dstType.FieldByName(fieldName)
		if !found {
			continue
		}
		foundValue := dstVal.FieldByName(fieldName)
		if foundValue.CanSet() {
			foundValue.Set(fieldValue)
		}
	}

	return nil
}
