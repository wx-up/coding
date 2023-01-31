//go:build e2e

package integration

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wx-up/coding/orm"
)

// SimpleStruct 测试模型
type SimpleStruct struct {
	Id   uint64
	Name string
}

type Suite struct {
	suite.Suite

	driver string
	dsn    string

	db *orm.DB
}

func (s *Suite) SetupSuite() {
	db, err := orm.Open(s.driver, s.dsn)
	require.NoError(s.T(), err)

	// 等待服务可用
	err = db.Wait()
	require.NoError(s.T(), err)

	s.db = db
}
