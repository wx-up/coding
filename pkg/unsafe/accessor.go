package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type Accessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

type FieldMeta struct {
	offset uintptr
	typ    reflect.Type
}

// NewAccessor 仅支持指向结构体的指针
func NewAccessor(entity any) (*Accessor, error) {
	if entity == nil {
		return nil, errors.New("entity 不能为 nil")
	}
	typ := reflect.TypeOf(entity)
	if !(typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct) {
		return nil, errors.New("不支持的类型")
	}
	res := &Accessor{
		fields: make(map[string]FieldMeta),
	}
	typ = typ.Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		// 填充字段的类型信息以及偏移量
		res.fields[field.Name] = FieldMeta{
			offset: field.Offset,
			typ:    field.Type,
		}
	}

	// 起始地址
	val := reflect.ValueOf(entity)
	res.entityAddr = val.UnsafePointer()
	return res, nil
}

// Field 获取字段的值
func (a *Accessor) Field(name string) (any, error) {
	if name == "" {
		return nil, errors.New("name 不能为空")
	}
	field, ok := a.fields[name]
	if !ok {
		return nil, errors.New("不存在的字段")
	}

	// 知道具体类型的操作：*(*int)(unsafe.Pointer(uintptr(u.entityAddr) + meta.offset))
	// 不知道具体类型的操作如下：（ 性能相对差一点 ）
	v := reflect.NewAt(field.typ, unsafe.Pointer(uintptr(a.entityAddr)+field.offset)).Elem()
	return v.Interface(), nil
}

// SetField 设置字段的值
func (a *Accessor) SetField(name string, val any) error {
	if name == "" {
		return errors.New("name 不能为空")
	}
	if val == nil {
		return errors.New("val 不能为nil")
	}
	meta, ok := a.fields[name]
	if !ok {
		return errors.New("不存在的字段")
	}
	refV := reflect.NewAt(meta.typ, unsafe.Pointer(uintptr(a.entityAddr)+meta.offset)).Elem()
	if refV.CanSet() {
		refV.Set(reflect.ValueOf(val))
	}
	return nil
}
