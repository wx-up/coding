package orm

import (
	"github.com/wx-up/coding/orm/internal/errs"
	"github.com/wx-up/coding/orm/internal/model"
)

/*
 migrate 的难点不在于 CREATE TABLE 建表语句，而在于 ALTER 语句

 ALTER 构建思路：
	解析结构体
	发起一次查询获取到表的列信息（ https://blog.csdn.net/Mint6/article/details/90321602 ）
	比较两者，再生成 ALTER 语句

还有一种实现思路：可以使用 sql 的抽象语法树（ https://github.com/xwb1989/sqlparser ）

一般公司内部，不采用 ORM 管理 DDL 语句，而是自己手写 DDL 语句，然后通过 DMS 审核执行（ 测试或者开发环境可以用 ORM 管理 ）
*/

// CreateTableBuilder 生成建表语句
type CreateTableBuilder struct {
	r   *model.Registry
	val any
}

func NewCreateTableBuilder() *CreateTableBuilder {
	return &CreateTableBuilder{}
}

func (b *CreateTableBuilder) Registry(r *model.Registry) *CreateTableBuilder {
	b.r = r
	return b
}

func (b *CreateTableBuilder) Val(val any) *CreateTableBuilder {
	b.val = val
	return b
}

/*
CREATE TABLE IS NOT EXIST (
)
*/

// Build 构建 create table 语句
func (b *CreateTableBuilder) Build() (*Query, error) {
	if b.r == nil {
		return nil, errs.NewErrParamEmpty("register")
	}
	if b.val == nil {
		return nil, errs.NewErrParamEmpty("val")
	}
	return nil, nil
}
