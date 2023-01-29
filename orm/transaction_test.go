package orm

import (
	"context"
	"errors"
	"testing"
)

/*
* 事务扩散方案
 */

type UserDao struct {
	sess Session
}

func (dao *UserDao) GetByIdV2(ctx context.Context, id int64) (obj *User, err error) {
	return NewSelector[User](dao.sess).Where(C("Id").Eq(id)).Get(ctx)
}

func TestUserDao(t *testing.T) {
	db, err := OpenDB(nil)
	if err != nil {
		t.Fatal(err)
	}
	dao1 := UserDao{
		sess: db,
	}
	_ = dao1

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	dao2 := UserDao{
		sess: tx,
	}
	_ = dao2
}

type User struct {
}

func (dao *UserDao) GetById(ctx context.Context, id int64) (obj *User, err error) {
	tx := ctx.Value("tx")
	if tx == nil {
		return nil, errors.New("未开启事务")
	} else {

	}
	return nil, nil
}
