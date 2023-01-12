package orm

import (
	"github.com/wx-up/coding/orm/internal/model"
	"strings"
)

type builder struct {
	// 元数据
	model *model.Model
	// SQL builder
	sb strings.Builder
	// SQL args
	args []any

	// 方言
	dialect Dialect
}

// addArgs 添加参数
func (b *builder) addArgs(args ...any) {
	b.args = append(b.args, args...)
}

// quote 添加引号
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.dialect.Quoter())
	b.sb.WriteString(name)
	b.sb.WriteByte(b.dialect.Quoter())
}
