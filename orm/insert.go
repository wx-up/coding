package orm

import (
	"context"
	"github.com/wx-up/coding/orm/internal/errs"
	"github.com/wx-up/coding/orm/internal/model"
)

type Inserter[T any] struct {
	sess   Session
	values []*T

	// 指定插入的列
	cs []string

	// 用于实现 upsert 语句，当出现冲突时更新指定的字段
	onConflict *ConflictKey

	builder
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	q, err := i.Build()
	if err != nil {
		return Result{err: err}
	}
	res, err := i.sess.execContext(ctx, q.SQL, q.Args...)
	return Result{
		err: err,
		res: res,
	}
}

type ConflictKey struct {
	assigns         []Assignable
	doNothing       bool
	conflictColumns []string
}

type ConflictKeyBuilder[T any] struct {
	i         *Inserter[T]
	assigns   []Assignable
	doNothing bool

	// 指定冲突的列，mysql 是不支持的，但是 sqlite3 是支持的
	conflictColumns []string
}

func (c *ConflictKeyBuilder[T]) Build() *ConflictKey {
	return &ConflictKey{
		assigns:         c.assigns,
		doNothing:       c.doNothing,
		conflictColumns: c.conflictColumns,
	}
}

// ConflictColumns 指定冲突的列
func (c *ConflictKeyBuilder[T]) ConflictColumns(cols []string) *ConflictKeyBuilder[T] {
	c.conflictColumns = cols
	return c
}

func (c *ConflictKeyBuilder[T]) Update(as ...Assignable) *Inserter[T] {
	c.assigns = as
	c.i.onConflict = c.Build()
	return c.i
}

func (c *ConflictKeyBuilder[T]) DoNothing() *Inserter[T] {
	c.doNothing = true
	c.i.onConflict = c.Build()
	return c.i
}

func (i *Inserter[T]) OnConflict() *ConflictKeyBuilder[T] {
	return &ConflictKeyBuilder[T]{
		i: i,
	}
}

// Update 这种实现方式没有上面这种 Builder 模式好，最起码调用 OnConflict 之后用户就很清楚后面的 Update 是在
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

func NewInserter[T any](sess Session) *Inserter[T] {
	return &Inserter[T]{
		builder: builder{
			dialect: sess.getCore().Dialect(),
		},
		sess: sess,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) <= 0 {
		return nil, errs.ErrInsertValuesEmpty
	}
	var err error
	t := new(T)
	i.model, err = i.sess.getCore().R().Get(t)
	if err != nil {
		return nil, err
	}
	i.sb.WriteString("INSERT INTO ")
	i.quote(i.model.TableName)
	i.sb.WriteByte(' ')
	i.sb.WriteByte('(')

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
			i.sb.WriteByte(',')
		}
		i.quote(col.ColName)
	}

	i.sb.WriteByte(')')
	i.sb.WriteString(" VALUES ")

	// map 和 slice 都是要预估容量的
	i.args = make([]any, 0, len(i.values)*len(columns))

	for valIndex, val := range i.values {
		if valIndex > 0 {
			i.sb.WriteByte(',')
		}
		valCreator := i.sess.getCore().ValCreator()(i.model, val)
		i.sb.WriteByte('(')
		for index, field := range columns {
			if index > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			arg, err := valCreator.Field(field.Name)
			if err != nil {
				return nil, err
			}
			//args = append(args, refVal.FieldByName(field.Name).Interface())
			i.addArgs(arg)
		}
		i.sb.WriteByte(')')
	}

	// 构建 ON DUPLICATE KEY UPDATE 语句
	if i.onConflict != nil {
		err = i.dialect.BuildOnConflict(&i.builder, i.onConflict)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}
