package valuer

import (
	"database/sql"
	"reflect"
	"unsafe"

	"github.com/wx-up/coding/orm/internal/errs"

	"github.com/wx-up/coding/orm/internal/model"
)

type unsafeValuer struct {
	addr  unsafe.Pointer
	model *model.Model
}

func (u unsafeValuer) Field(name string) (any, error) {
	field, ok := u.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	ptr := unsafe.Pointer(uintptr(u.addr) + field.Offset)
	// 确切类型的方式读取：*(*int)(ptr)
	// 下面是任意类型的读取：
	val := reflect.NewAt(field.Typ, ptr)
	return val, nil
}

var _ Factory = NewUnsafeValuer

func NewUnsafeValuer(model *model.Model, t any) Valuer {
	return &unsafeValuer{
		addr:  reflect.ValueOf(t).UnsafePointer(),
		model: model,
	}
}

func (u unsafeValuer) SetColumns(rows *sql.Rows) error {
	// 获取数据库返回的列名
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 当查询的列多于 model 中定义的列时，则报错
	if len(columns) > len(u.model.ColumnMap) {
		return errs.ErrTooManyColumns
	}

	// 将数据库返回的值 scan 到 vs 中
	vs := make([]any, 0, len(columns))
	for _, col := range columns {
		// 当查询的列在模型中未定义时则报错
		fd, ok := u.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}

		// 在指定位置创建变量，reflect.New 在任意位置创建变量
		// 字段真实的地址 = 结构体的起始地址 + 字段的偏移量
		reflectV := reflect.NewAt(fd.Typ, unsafe.Pointer(uintptr(u.addr)+fd.Offset))
		vs = append(vs, reflectV.Interface())
	}

	return rows.Scan(vs...)
}
