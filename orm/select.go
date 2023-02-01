package orm

import (
	"context"
	"fmt"

	"github.com/wx-up/coding/orm/internal/errs"
)

type Selector[T any] struct {
	builder

	// 表名
	tbl TableReference

	// WHERE 条件
	ps []Predicate

	columns []Selectable

	sess Session

	// 中间件
	ms []Middleware
}

func NewSelector[T any](sess Session) *Selector[T] {
	core := sess.getCore()
	s := &Selector[T]{
		sess: sess,
		// 获取设置在 DB 对象上的 middlewares
		ms: core.Middlewares(),
	}
	s.dialect = core.Dialect()
	return s
}

// Use 添加 Selector 层面的 middleware
func (s *Selector[T]) Use(ms []Middleware) *Selector[T] {
	s.ms = append(s.ms, ms...)
	return s
}

type Selectable interface {
	selectable()
}

/*
// 给标记接口具体的实现
type Selectable interface {
	selectable(b *builder) error
}
*/

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.ps = ps
	return s
}

func (s *Selector[T]) From(tbl TableReference) *Selector[T] {
	s.tbl = tbl
	return s
}

func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	return s
}

func (s *Selector[T]) OrderBy(order ...OrderBy) *Selector[T] {
	return s
}

// TableName 返回表名
// 如果用户传递了表名就直接按照用户传递的，不做检测它是否携带反引号，没有则追加的逻辑（ 这只是设计上的决策没有对错 ）
// 如果用户没有传递表名则使用结构体名称的复数，并前后添加反引号
//func (s *Selector[T]) TableName() string {
//	tblName := s.tbl
//	if tblName != "" {
//		return tblName
//	}
//	// 结构体名称的复数
//	return "`" + s.model.TableName + "`"
//}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	t := new(T)
	s.model, err = s.sess.getCore().R().Get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")

	err = s.buildTable(s.tbl)
	if err != nil {
		return nil, err
	}

	// 拼接 where
	if len(s.ps) > 0 {
		// WHERE 是可选的，当它存在的时候需要前后加上空格
		s.sb.WriteString(" WHERE ")

		// 多个 predicate 之间使用 and 连接
		pre := s.ps[0]
		for i := 1; i < len(s.ps); i++ {
			pre = pre.And(s.ps[i])
		}

		if err := s.buildExpression(pre); err != nil {
			return nil, err
		}
	}

	// 拼接分号
	s.sb.WriteByte(';')

	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildTable(tbl TableReference) error {
	switch tr := tbl.(type) {
	case nil: // 不调用 From 函数
		s.quote(s.model.TableName)
	case Table:
		model, err := s.sess.getCore().R().Get(tr.entry)
		if err != nil {
			return err
		}
		s.quote(model.TableName)
	// TableOf(&User{}).Join(TableOf(&User{}).Join(&Order{}))
	case Join:
		s.sb.WriteByte('(')
		err := s.buildTable(tr.left)
		if err != nil {
			return err
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(tr.typ)
		s.sb.WriteByte(' ')

		err = s.buildTable(tr.right)
		if err != nil {
			return nil
		}

		// 处理 ON，逻辑和 where 一样
		if len(tr.on) > 0 {
			s.sb.WriteString(" ON ")
			pre := s.ps[0]
			for i := 1; i < len(s.ps); i++ {
				pre = pre.And(s.ps[i])
			}

			if err = s.buildExpression(pre); err != nil {
				return err
			}
		}

		// 处理 Using
		if len(tr.using) > 0 {
			s.sb.WriteString(" USING(")
			for _, col := range tr.using {
				_ = col
			}
		}
		s.sb.WriteByte(')')
	case SubQuery:
	default:
		return errs.NewErrUnsupportedTableReference(tbl)
	}
	return nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteString("*")
		return nil
	}
	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch c := col.(type) {
		case Column: // 列
			if err := s.buildColumn(c); err != nil {
				return err
			}
		case Aggregate: // 聚合函数
			fd, ok := s.model.FieldMap[c.arg]
			if !ok {
				return errs.NewErrUnknownField(c.arg)
			}
			s.sb.WriteString(c.fn)
			s.sb.WriteByte('(')
			s.sb.WriteByte('`')
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte('`')
			s.sb.WriteByte(')')
			if c.alias != "" {
				s.sb.WriteString(" AS ")
				s.sb.WriteString(c.alias)
			}
		case RawExpr:
			// SELECT xxx  其中 xxx 可以是很复杂的表达式，比如函数调用等等
			// 所以要预留 args 字段，作为参数
			s.sb.WriteString(c.raw)
			if len(c.args) > 0 {
				s.args = append(s.args, c.args...)
			}
		}
	}
	return nil
}

func (s *Selector[T]) buildColumn(c Column) error {
	// colName = `user`.`id`
	colName, err := s.column(s.tbl, c)
	if err != nil {
		return err
	}

	s.sb.WriteString(colName)

	// 处理别名
	if c.alias != "" {
		s.sb.WriteString(" AS ")
		s.sb.WriteString(c.alias)
	}
	return nil
}

func (s *Selector[T]) column(tr TableReference, c Column) (string, error) {
	// 断言 TableReference 类型
	switch tbl := tr.(type) {
	case nil:
		fd, ok := s.model.FieldMap[c.name]
		if !ok {
			return "", errs.NewErrUnknownField(c.name)
		}

		// 添加表名前缀，当 join 查询时，字段存在在多张表时，没有表名限定的话，会出错
		s.quote(s.model.TableName)
		s.sb.WriteByte('.')

		return fmt.Sprintf("`%s`.`%s`", s.model.TableName, fd.ColName), nil

	case Table:
		model, err := s.sess.getCore().R().Get(tbl.entry)
		if err != nil {
			return "", err
		}
		fd, ok := model.FieldMap[c.name]
		if !ok {
			return "", errs.NewErrUnknownField(c.name)
		}
		return fmt.Sprintf("`%s`.`%s`", model.TableName, fd.ColName), nil

	case Join:
		// 递归处理
		colName, err := s.column(tbl.left, c)
		if err == nil {
			return colName, nil
		}
		// 如果左边不在，那就找右边
		return s.column(tbl.right, c)
	default:
		return "", errs.NewErrUnsupportedTableReference(s.tbl)
	}
}

// Get 查询单条数据
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	// 生成 SQL 以及 参数
	query, err := s.Build()
	if err != nil {
		return nil, err
	}

	// 最后也是委托给 原生查询
	return RawQuery[T](s.sess, query.SQL, query.Args...).Use(s.ms).Get(ctx)
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Selector[T]) buildExpression(expr expression) error {
	if expr == nil {
		return nil
	}
	switch v := expr.(type) {
	case RawExpr:
		s.sb.WriteString(v.raw)
		if s.args == nil {
			s.args = make([]any, 0, len(v.args))
		}
		s.args = append(s.args, v.args...)
	case Column:
		s.sb.WriteByte('`')

		// 检测用户传递的字段 ( C("Name") ) 是否合法
		fd, ok := s.model.FieldMap[v.name]
		if !ok {
			return errs.NewErrUnknownField(v.name)
		}

		s.sb.WriteString(fd.ColName)
		s.sb.WriteByte('`')
	case Val:
		s.sb.WriteByte('?')

		// 预估容量
		if s.args == nil {
			s.args = make([]any, 0, 4)
		}
		s.args = append(s.args, v.val)

	case Predicate:
		// 如果左表达式是 predicate 则添加括号
		_, ok := v.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(v.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

		// 当操作符不为空的时候再处理
		if v.op != opEmpty {
			// 操作符
			if v.op != opNot {
				s.sb.WriteString(" ")
			}
			s.sb.WriteString(v.op.String())
			s.sb.WriteString(" ")
		}

		// 如果右表达式是 predicate 则添加括号
		_, ok = v.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(v.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

	default:
		return errs.NewErrUnsupportedExpressionType(expr)
	}
	return nil
}
