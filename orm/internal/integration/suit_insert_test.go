package integration

import (
	"github.com/stretchr/testify/suite"
	"github.com/wx-up/coding/orm"
)

type InsertTestSuit struct {
	suite.Suite
	db *orm.DB
}

func (suit *InsertTestSuit) TestMysql() {
	t := suit.T()
}
