package orm

import (
	"context"
	"github.com/wx-up/coding/orm/internal/errs"
)

type Selector[T any] struct {
	builder

	// 表名
	tbl string

	// WHERE 条件
	ps []Predicate

	db *DB

	columns []Selectable
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

type Selectable interface {
	selectable()
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.ps = ps
	return s
}

func (s *Selector[T]) From(tbl string) *Selector[T] {
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
	s.model, err = s.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")

	if s.tbl != "" {
		s.sb.WriteString(s.tbl)
	} else {
		s.sb.WriteString(s.model.TableName)
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
			fd, ok := s.model.FieldMap[c.name]
			if !ok {
				return errs.NewErrUnknownField(c.name)
			}
			s.sb.WriteByte('`')
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte('`')
			// 如果别名存在，则设置
			if c.alias != "" {
				s.sb.WriteString(" AS ")
				s.sb.WriteString(c.alias)
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

// Get 查询单条数据
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	query, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := s.db.db.QueryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, ErrNoRows
	}

	t := new(T)

	return t, s.db.valCreator(s.model, t).SetColumns(rows)
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
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
