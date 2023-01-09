package valuer

import (
	"database/sql"
	"github.com/wx-up/coding/orm/internal/model"
)

/*
 从实现上可以看出：

	反射的逻辑是先根据 reflect.Type 创建对象并通过 Scan 赋值，然后再通过反射的 Set 设置结构体的值
	因为 reflect.New 创建的对象并不知道在哪里，所以会有 for Set 的操作，将值拷贝到指定的位置

	unsafe 基于 reflect.NewAt 直接在指定位置创建对象，Scan 即可，不需要反射 挪到指定位置 这一步

*/

type Valuer interface {
	SetColumns(rows *sql.Rows) error
}

type Factory func(model *model.Model, t any) Valuer
