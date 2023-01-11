package orm

import (
	"github.com/wx-up/coding/orm/internal/errs"
	"github.com/wx-up/coding/orm/internal/model"
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	db     *DB
	model  *model.Model
	values []*T

	// 指定插入的列
	cs []string

	// 用于实现 upsert 语句，当出现冲突时更新指定的字段
	onConflict *ConflictKey
}

type ConflictKey struct {
	assigns []Assignable
}

type ConflictKeyBuilder[T any] struct {
	i       *Inserter[T]
	assigns []Assignable
}

func (c *ConflictKeyBuilder[T]) Build() *ConflictKey {
	return &ConflictKey{assigns: c.assigns}
}

func (c *ConflictKeyBuilder[T]) Update(as ...Assignable) *Inserter[T] {
	c.assigns = as
	c.i.onConflict = c.Build()
	return c.i
}

func (i *Inserter[T]) OnConflict() *ConflictKeyBuilder[T] {
	return &ConflictKeyBuilder[T]{
		i: i,
	}
}

// 这种实现方式没有上面这种 Builder 模式好，最起码调用 OnConflict 之后用户就很清楚后面的 Update 是在
// 冲突时指定更新的列
// 如果这种实现方式，直接 Update 就不够直观
//func (i *Inserter[T]) Update(as ...Assignable) *Inserter[T] {
//	i.onConflict = &ConflictKey{
//		assigns: as,
//	}
//	return i
//}

// Columns 也可以按照 SELECT 的 Columns 的设计 -- 结构化
func (i *Inserter[T]) Columns(cs ...string) *Inserter[T] {
	i.cs = cs
	return i
}

// Values 指定要插入的数据，结构体指针
func (i *Inserter[T]) Values(vs ...*T) *Inserter[T] {
	i.values = vs
	return i
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) <= 0 {
		return nil, errs.ErrInsertValuesEmpty
	}
	var (
		err     error
		builder strings.Builder
	)
	t := new(T)
	i.model, err = i.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	builder.WriteString("INSERT INTO ")
	builder.WriteByte('`')
	builder.WriteString(i.model.TableName)
	builder.WriteByte('`')
	builder.WriteByte(' ')
	builder.WriteByte('(')

	// 不指定列，则为全部的列
	columns := i.model.Columns
	if length := len(i.cs); length > 0 { // 指定的列
		columns = make([]*model.Field, 0, length)
		for _, col := range i.cs {
			fd, ok := i.model.FieldMap[col]
			if !ok {
				return nil, errs.NewErrUnknownField(col)
			}
			columns = append(columns, fd)
		}
	}

	for index, col := range columns {
		if index > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('`')
		builder.WriteString(col.ColName)
		builder.WriteByte('`')
	}

	builder.WriteByte(')')
	builder.WriteString(" VALUES ")

	// map 和 slice 都是要预估容量的
	args := make([]any, 0, len(i.values)*len(i.model.Columns))

	for valIndex, val := range i.values {
		if valIndex > 0 {
			builder.WriteByte(',')
		}

		// 如果是指针，则转成结构体
		refVal := reflect.ValueOf(val)
		for refVal.Kind() == reflect.Ptr {
			refVal = refVal.Elem()
		}

		builder.WriteByte('(')
		for index, field := range columns {
			if index > 0 {
				builder.WriteByte(',')
			}
			builder.WriteByte('?')
			//args = append(args, refVal.FieldByName(field.Name).Interface())
			args = append(args, refVal.FieldByIndex(field.Index).Interface())
		}
		builder.WriteByte(')')
	}

	// 构建 ON DUPLICATE KEY UPDATE 语句
	if i.onConflict != nil {
		builder.WriteString(" ON DUPLICATE KEY UPDATE ")
		for index, assign := range i.onConflict.assigns {
			if index > 0 {
				builder.WriteByte(',')
			}
			switch v := assign.(type) {
			case Assigment: // 用于生成 name=?
				fd, ok := i.model.FieldMap[v.column]
				if !ok {
					return nil, errs.NewErrUnknownField(v.column)
				}
				builder.WriteByte('`')
				builder.WriteString(fd.ColName)
				builder.WriteByte('`')
				builder.WriteString("=?")
				args = append(args, v.val)
			case Column: // 用于生成 name=VALUES(name)
				fd, ok := i.model.FieldMap[v.name]
				if !ok {
					return nil, errs.NewErrUnknownField(v.name)
				}
				builder.WriteByte('`')
				builder.WriteString(fd.ColName)
				builder.WriteByte('`')
				builder.WriteString("=VALUES(")
				builder.WriteByte('`')
				builder.WriteString(fd.ColName)
				builder.WriteByte('`')
				builder.WriteByte(')')
			}
		}
	}

	builder.WriteByte(';')
	return &Query{
		SQL:  builder.String(),
		Args: args,
	}, nil
}
