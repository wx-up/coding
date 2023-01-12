package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wx-up/coding/orm/internal/errs"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	type TestModel struct {
		Id        int64
		FirstName string
		Age       int64
		LastName  *sql.NullString
	}
	db, err := OpenDB(nil)
	require.Nil(t, err)
	tests := []struct {
		name     string
		inserter QueryBuilder
		want     *Query
		wantErr  error
	}{
		{
			name:     "no row", // 参数错误
			inserter: NewInserter[TestModel](db).Values(),
			wantErr:  errs.ErrInsertValuesEmpty,
		},
		{
			name: "insert one row", // 插入一行记录，全部字段
			inserter: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "哈哈",
				Age:       24,
				LastName: &sql.NullString{
					String: "嘻嘻",
					Valid:  true,
				},
			}),
			want: &Query{
				SQL:  "INSERT INTO `test_models` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?);",
				Args: []any{int64(1), "哈哈", int64(24), &sql.NullString{String: "嘻嘻", Valid: true}},
			},
		},
		{
			name: "insert many row", // 插入多行，全部字段
			inserter: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "哈哈",
					Age:       24,
					LastName: &sql.NullString{
						String: "嘻嘻",
						Valid:  true,
					},
				}, &TestModel{
					Id:        1,
					FirstName: "哈哈",
					Age:       24,
					LastName: &sql.NullString{
						String: "嘻嘻",
						Valid:  true,
					},
				}),
			want: &Query{
				SQL: "INSERT INTO `test_models` (`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{
					int64(1), "哈哈", int64(24), &sql.NullString{String: "嘻嘻", Valid: true},
					int64(1), "哈哈", int64(24), &sql.NullString{String: "嘻嘻", Valid: true},
				},
			},
		},
		{
			name: "specify column", // 指定列
			inserter: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "哈哈",
				Age:       24,
				LastName: &sql.NullString{
					String: "嘻嘻",
					Valid:  true,
				},
			}).Columns("FirstName", "Age"),
			want: &Query{
				SQL:  "INSERT INTO `test_models` (`first_name`,`age`) VALUES (?,?);",
				Args: []any{"哈哈", int64(24)},
			},
		},
		{
			name: "upsert", // upsert 语句构建
			inserter: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "哈哈",
				Age:       24,
				LastName: &sql.NullString{
					String: "嘻嘻",
					Valid:  true,
				},
			}).Columns("FirstName", "Age").OnConflict().Update(Assign("Age", 100)),
			want: &Query{
				SQL:  "INSERT INTO `test_models` (`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{"哈哈", int64(24), 100},
			},
		},
		{
			name: "upsert use column", // upsert 语句构建
			inserter: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "哈哈",
				Age:       24,
				LastName: &sql.NullString{
					String: "嘻嘻",
					Valid:  true,
				},
			}).Columns("FirstName", "Age").OnConflict().Update(C("Age"), Assign("FirstName", "Bob")),
			want: &Query{
				// 出现冲突的时候，则更新，但是 age 的值不变，first_name 设置为 Bob
				SQL:  "INSERT INTO `test_models` (`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `age`=VALUES(`age`),`first_name`=?;",
				Args: []any{"哈哈", int64(24), "Bob"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.inserter.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
