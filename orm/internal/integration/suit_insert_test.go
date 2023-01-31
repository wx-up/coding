package integration

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wx-up/coding/orm"
	"testing"
)

type InsertTestSuit struct {
	Suite
}

func (s *InsertTestSuit) TestInsert() {
	testCases := []struct {
		name string
		i    *orm.Inserter[SimpleStruct]

		wantData     *SimpleStruct
		wantErr      error
		rowsAffected int64
	}{
		{
			name: "单行",
			i: orm.NewInserter[SimpleStruct](s.db).Values(&SimpleStruct{
				Name: "单行",
			}),

			// 涉及到数据存储的外部环境时，在测试完成之后，需要去对比数据是否和期望的一致
			// 在当前场景下：插入数据之后，查询一次，查看数据库的记录是否和期望的一致
			wantData: &SimpleStruct{
				Id:   1,
				Name: "单行",
			},
			wantErr:      nil,
			rowsAffected: 1,
		},
	}

	t := s.T()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.rowsAffected, affected)

			id, _ := res.LastInsertId()
			obj, _ := orm.NewSelector[SimpleStruct](s.db).Where(orm.C("Id").Eq(id)).Get(context.Background())
			assert.Equal(t, tc.wantData, obj)
		})
	}
}

// TearDownTest 每个测试结束都会运行的钩子
// 保证测试之间相互独立（ 删除自己的测试数据 ）
func (s *InsertTestSuit) TearDownTest() {
	res := orm.RawQuery[any](s.db, "TRUNCATE TABLE `simple_structs`;").Exec(context.Background())
	require.NoError(s.T(), res.Err())
}

func Test_MYSQL8_Insert(t *testing.T) {
	suite.Run(t, &InsertTestSuit{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}
