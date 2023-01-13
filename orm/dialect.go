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
	if len(onConflict.assigns) <= 0 {
		return errs.ErrAssigmentEmpty
	}
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for index, assign := range onConflict.assigns {
		if index > 0 {
			b.sb.WriteByte(',')
		}

		// 其实可以将标记接口实现，类似 gorm 的做法，但是会有一个难点：在 SQL 中一个实体放在不同部分的时候，构建逻辑是不同的
		// 比如 aggregate 实体 放在 SELECT 中是可以有别名的：
		// 	SELECT AVG(`age`) as avg_age
		// 放在 HAVING 中则是不能用别名的
		//  HAVING AVG(`age`) < 10
		// switch 的好处就是逻辑清晰，可读性强，当 case 比较多的时候维护性也比较差
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
	if len(onConflict.conflictColumns) <= 0 {
		return errs.ErrConflictEmpty
	}

	b.sb.WriteString(" ON CONFLICT ")
	b.sb.WriteByte('(')
	for index, col := range onConflict.conflictColumns {
		if index > 0 {
			b.sb.WriteByte(',')
		}
		fd, ok := b.model.FieldMap[col]
		if !ok {
			return errs.NewErrUnknownField(col)
		}
		b.quote(fd.ColName)
	}
	b.sb.WriteByte(')')

	switch onConflict.doNothing {
	case false:
		if len(onConflict.assigns) <= 0 {
			return errs.ErrAssigmentEmpty
		}
		b.sb.WriteString(" DO UPDATE SET ")
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
			case Column: // 用于生成 column=excluded.column （ excluded 为 sqlite3 的关键字类似 VALUES ）
				fd, ok := b.model.FieldMap[v.name]
				if !ok {
					return errs.NewErrUnknownField(v.name)
				}
				b.quote(fd.ColName)
				b.sb.WriteString("=excluded.")
				b.quote(fd.ColName)
			default:
				return errs.NewErrUnsupportedAssignableType(assign)
			}
		}
	case true:
		b.sb.WriteString(" NOTHING")
	}

	return nil
}
