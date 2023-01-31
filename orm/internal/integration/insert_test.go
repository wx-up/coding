//go:build e2e

// 集成测试通常是依赖于外部环境，以及数据准备，而且跑的比较慢，所以一般不会和单元测试集成到一起
// 因此会给集成测试打一个 tag 比如：e2e（ 端到端 ）

package integration

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wx-up/coding/orm"
	"testing"
)

// TestMysql_Insert 集成测试
func TestMysql_Insert(t *testing.T) {
	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	require.NoError(t, err)

	// 存在容器启动了，但是 mysql 服务还没有就绪，所以这里需要 wait 一下等待 mysql 就绪
	err = db.Wait()
	require.NoError(t, err)

	// 结束之后，删除测试数据
	defer func() {
		orm.RawQuery[any](db, "TRUNCATE TABLE `simple_structs`;").Exec(context.Background())
	}()

	testCases := []struct {
		name            string
		i               *orm.Inserter[SimpleStruct]
		wantRowAffected int64
		wantErr         error

		wantData *SimpleStruct
	}{
		{
			name: "single insert",
			i: orm.NewInserter[SimpleStruct](db).Values(&SimpleStruct{
				Name: "single insert",
			}),
			wantRowAffected: 1,
			wantData:        &SimpleStruct{Id: 1, Name: "single insert"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			rows, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRowAffected, rows)
			id, _ := res.LastInsertId()
			obj, err := orm.NewSelector[SimpleStruct](db).Where(orm.C("Id").Eq(id)).Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, obj)
		})
	}
}
