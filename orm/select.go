package orm

import (
	"context"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	// SQL builder
	sb strings.Builder

	// SQL args
	args []any

	// 表名
	tbl string

	// WHERE 条件
	ps []Predicate
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.ps = ps
	return s
}

func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.tbl = tbl
	return s
}

// TableName 返回表名
// 如果用户传递了表名就直接按照用户传递的，不做检测它是否携带反引号，没有则追加的逻辑（ 这只是设计上的决策没有对错 ）
// 如果用户没有传递表名则使用结构体名称的复数，并前后添加反引号
func (s *Selector[T]) TableName() string {
	tblName := s.tbl
	if tblName != "" {
		return tblName
	}
	// 结构体名称的复数
	var t T
	typ := reflect.TypeOf(t)
	tblName = strcase.ToSnake(pluralize.NewClient().Plural(typ.Name()))
	return "`" + tblName + "`"
}

func (s *Selector[T]) Build() (*Query, error) {

	s.sb.WriteString("SELECT * FROM ")
	s.sb.WriteString(s.TableName())

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

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
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
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(v.name)
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

		// 操作符
		if v.op != opNot {
			s.sb.WriteString(" ")
		}
		s.sb.WriteString(v.op.String())
		s.sb.WriteString(" ")

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
