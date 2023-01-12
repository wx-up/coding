package orm

import (
	"github.com/wx-up/coding/orm/internal/errs"
)

// Dialect 方言
type Dialect interface {
	// Quoter 返回一个引号，引用列名、表名的引号
	// mysql 是反引号
	// pgsql 是双引号
	Quoter() byte

	// BuildOnConflict upsert 语句
	BuildOnConflict(b *builder, key *ConflictKey) error
}

// standardSQL 标准SQL的实现
type standardSQL struct {
}

// MysqlDialect mysql 方言的实现
type MysqlDialect struct {
	standardSQL
}

func NewMysqlDialect() *MysqlDialect {
	return &MysqlDialect{}
}

func (m *MysqlDialect) Quoter() byte {
	return '`'
}

// BuildOnConflict 内部实现需要 Inserter 的 model 和 args  支持
// 在 Selector 中也有类似的需求，所以将它们进行抽象，得到 builder 结构体
func (m *MysqlDialect) BuildOnConflict(b *builder, onConflict *ConflictKey) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for index, assign := range onConflict.assigns {
		if index > 0 {
			b.sb.WriteByte(',')
		}
		switch v := assign.(type) {
		case Assigment: // 用于生成 name=?
			fd, ok := b.model.FieldMap[v.column]
			if !ok {
				return errs.NewErrUnknownField(v.column)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.addArgs(v.val)
		case Column: // 用于生成 name=VALUES(name)
			fd, ok := b.model.FieldMap[v.name]
			if !ok {
				return errs.NewErrUnknownField(v.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sb.WriteByte(')')
		default:
			return errs.NewErrUnsupportedAssignableType(assign)
		}
	}
	return nil

}

// SqliteDialect sqlite 方言的实现
type SqliteDialect struct {
	standardSQL
}

func (s *SqliteDialect) Quoter() byte {
	return '`'
}

func (s *SqliteDialect) BuildOnConflict(b *builder, onConflict *ConflictKey) error {
	return nil
}
