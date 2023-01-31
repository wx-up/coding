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

func Test_MYSQL8_Select(t *testing.T) {
	suite.Run(t, &SelectTestSuite{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

type SelectTestSuite struct {
	Suite
}

func (s *SelectTestSuite) SetupSuite() {
	s.Suite.SetupSuite()

	// 准备数据
	res := orm.NewInserter[SimpleStruct](s.db).Values(
		&SimpleStruct{Name: "测试一"},
		&SimpleStruct{Name: "测试二"},
		&SimpleStruct{Name: "测试三"},
	).Exec(context.Background())

	require.NoError(s.T(), res.Err())
}

func (s *SelectTestSuite) TearDownSuite() {
	res := orm.RawQuery[SimpleStruct](s.db, "TRUNCATE TABLE `simple_structs`;").Exec(context.Background())
	require.NoError(s.T(), res.Err())
}

func (s *SelectTestSuite) TestSelect() {
	testCases := []struct {
		name string

		s *orm.Selector[SimpleStruct]

		wantErr  error
		wantData any
	}{
		{
			name:    "未找到",
			s:       orm.NewSelector[SimpleStruct](s.db).Where(orm.C("Name").Eq("测试十")),
			wantErr: orm.ErrNoRows,
		},
		{
			name:    "找到一行",
			s:       orm.NewSelector[SimpleStruct](s.db).Where(orm.C("Name").Eq("测试一")),
			wantErr: nil,
			wantData: &SimpleStruct{
				Id:   1,
				Name: "测试一",
			},
		},
	}
	t := s.T()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantData, obj)
		})
	}
}
