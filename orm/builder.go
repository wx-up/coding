package orm

import (
	"strings"

	"github.com/wx-up/coding/orm/internal/model"
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
	if b.args == nil {
		// WHERE 很少有查询能超过八个参数的
		// INSERT 是在外部初始化的，i.args = make([]any, 0, len(i.values)*len(columns))
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
}

// quote 添加引号
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.dialect.Quoter())
	b.sb.WriteString(name)
	b.sb.WriteByte(b.dialect.Quoter())
}
