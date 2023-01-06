package orm

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
)

type Model struct {
	// 结构体对应的表名
	tableName string
	// 结构体所有字段的元数据信息
	fieldMap map[string]*field
}

type field struct {
	// 当前字段对应的列名
	colName string
}

// parseModel 按照目前的设计其实不需要 error
// 但是根据经验随着功能迭代后续会需要 error，那时候再修改调用方会比较麻烦，所以现在就直接添加上
func parseModel(val any) (*Model, error) {
	if val == nil {
		return nil, errs.ErrParseModelValType
	}

	typ := reflect.TypeOf(val)

	// 处理多级指针的情况
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrParseModelValType
	}
	fieldCnt := typ.NumField()
	fieldMap := make(map[string]*field, fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		f := typ.Field(i)
		fieldMap[f.Name] = &field{
			colName: strcase.ToSnake(f.Name),
		}
	}
	return &Model{
		tableName: strcase.ToSnake(pluralize.NewClient().Plural(typ.Name())),
		fieldMap:  fieldMap,
	}, nil
}
