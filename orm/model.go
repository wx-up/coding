package orm

import (
	"fmt"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
)

// ModelOpt 带校验的 option 模式，有问题则返回 error
type ModelOpt func(*Model) error

func ModelWithTableName(tblName string) ModelOpt {
	return func(model *Model) error {
		model.tableName = tblName
		return nil
	}
}

// ModelWithColumnName 修改字段的列名
// 这种设计不够通用，当 field 结构体中新增了其他字段比如 autoincrement
// 那么又要写一个类似这样子的方法 ModelWithColumnAutoincrement
// 方法膨胀的比较快，相比之下 ModelWithColumn 更加通用
func ModelWithColumnName(field string, column string) ModelOpt {
	return func(model *Model) error {
		fd, ok := model.fieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = column
		return nil
	}
}

// ModelWithColumn 直接让用户传递一个 field 结构体，更加通用
func ModelWithColumn(field string, col *field) ModelOpt {
	return func(model *Model) error {
		_, ok := model.fieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		model.fieldMap[field] = col
		return nil
	}
}

type Model struct {
	// 结构体对应的表名
	tableName string
	// 结构体所有字段的元数据信息
	fieldMap map[string]*field

	// 数据库列名的映射
	columnMap map[string]*field
}

type field struct {
	// 当前字段名
	name string
	// 当前字段对应的列名
	colName string
	// 当前字段的类型，用于结果集生成
	typ reflect.Type
}

// TableName 自定义表名
// 用户可以通过实现 TableName 结构来返回不同的表名，从而实现分表的逻辑
type TableName interface {
	TableName() string
}

// DBName 自定义库
type DBName interface {
	DBName() string
}

/*
 orm 简单实现分库分表
*/

// Order orm 支持分库分表的话，在 Build SQL 的时候需要动态调用 DBName 和 TableName 而不是直接取 parseModel 的得到的值
// 否则库名和表名都固定了
type Order struct {
	Id     int64
	Region string
}

// DBName 分库
func (o *Order) DBName() string {
	return fmt.Sprintf("%s_order_db_%4d", o.Region, o.Id%1000)
}

// TableName 分表
func (o *Order) TableName() string {
	return fmt.Sprintf("order_%4d", o.Id%1000)
}
