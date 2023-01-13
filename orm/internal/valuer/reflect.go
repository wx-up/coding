package valuer

import (
	"database/sql"
	"github.com/wx-up/coding/orm/internal/errs"
	"github.com/wx-up/coding/orm/internal/model"
	"reflect"
)

type reflectValuer struct {
	model *model.Model
	// 泛型 T 的指针类型
	t any
}

func (r reflectValuer) Field(name string) (any, error) {
	// 如果是指针，则转成结构体
	refVal := reflect.ValueOf(r.t)
	for refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}
	field, ok := r.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}

	return refVal.FieldByIndex(field.Index).Interface(), nil
}

var _ Factory = NewReflectValuer

func NewReflectValuer(model *model.Model, t any) Valuer {
	return &reflectValuer{
		model: model,
		t:     t,
	}
}

func (r reflectValuer) SetColumns(rows *sql.Rows) error {
	// 获取数据库返回的列名
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 当查询的列多于 model 中定义的列时，则报错
	if len(columns) > len(r.model.ColumnMap) {
		return errs.ErrTooManyColumns
	}

	// 缓存，避免多次反射
	reflectVs := make([]reflect.Value, 0, len(columns))

	// 将数据库返回的值 scan 到 vs 中
	vs := make([]any, 0, len(columns))
	for _, col := range columns {
		// 当查询的列在模型中未定义时则报错
		fd, ok := r.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}

		reflectV := reflect.New(fd.Typ)
		vs = append(vs, reflectV.Interface())
		reflectVs = append(reflectVs, reflectV)
	}

	err = rows.Scan(vs...)
	if err != nil {
		return err
	}

	value := reflect.ValueOf(r.t).Elem()
	for index, col := range columns {
		// 根据列名获取字段名
		fd := r.model.ColumnMap[col]

		vField := value.FieldByName(fd.Name)
		if !vField.CanSet() {
			continue
		}
		// 通过反射设置字段的值
		vField.Set(reflectVs[index].Elem())
	}
	return nil
}
