package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wx-up/coding/orm/internal/errs"
	"github.com/wx-up/coding/orm/internal/valuer"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int
	LastName  *sql.NullString
}

func (TestModel) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS test_models(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
}

func TestSelect_Build(t *testing.T) {
	db, err := OpenDB(nil)
	require.Nil(t, err)
	tests := []struct {
		name    string
		builder QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "from",
			builder: NewSelector[TestModel](db).From("test_model_tbl"),
			want: &Query{
				SQL: "SELECT * FROM test_model_tbl;",
			},
			wantErr: nil,
		},
		{
			name:    "from with quotation mark", // 反引号
			builder: NewSelector[TestModel](db).From("`test_model_tbl`"),
			want: &Query{
				SQL: "SELECT * FROM `test_model_tbl`;",
			},
			wantErr: nil,
		},
		{
			name:    "no from",
			builder: NewSelector[TestModel](db),
			want: &Query{
				SQL: "SELECT * FROM `test_models`;",
			},
			wantErr: nil,
		},
		{
			name:    "from but empty",
			builder: NewSelector[TestModel](db).From(""),
			want: &Query{
				SQL: "SELECT * FROM `test_models`;",
			},
			wantErr: nil,
		},
		{
			name:    "from with db",
			builder: NewSelector[TestModel](db).From("byn.test_model"),
			want: &Query{
				SQL: "SELECT * FROM byn.test_model;",
			},
			wantErr: nil,
		},
		{
			name:    "single predicate",
			builder: NewSelector[TestModel](db).Where(C("Id").Eq(12)),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE `id` = ?;",
				Args: []any{12},
			},
			wantErr: nil,
		},
		{
			name:    "and predicates",
			builder: NewSelector[TestModel](db).Where(C("Age").Gt(12).And(C("Age").Lt(24))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{12, 24},
			},
			wantErr: nil,
		},
		{
			name:    "or predicates",
			builder: NewSelector[TestModel](db).Where(C("FirstName").Eq("bob").Or(C("Age").Gt(12))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`first_name` = ?) OR (`age` > ?);",
				Args: []any{"bob", 12},
			},
		},
		{
			name:    "not predicates",
			builder: NewSelector[TestModel](db).Where(C("FirstName").Eq("bob").And(Not(C("Age").Eq(12)))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`first_name` = ?) AND (NOT (`age` = ?));",
				Args: []any{"bob", 12},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithUnsafeValCreator())
	require.NoError(t, err)
	tests := []struct {
		name     string
		query    string
		mockErr  error         // mock 错误
		mockRows *sqlmock.Rows // mock 返回值
		wantErr  error
		wantVal  any
	}{
		// 这个测试用例主要是测试，当发生错误时是否返回错误
		{
			name:    "query invalid",
			query:   "SELECT .*",
			mockErr: errors.New("invalid query"),
			wantErr: errors.New("invalid query"),
		},
		// 没有数据时返回错误
		{
			name:     "no row",
			query:    "SELECT .*",
			mockErr:  nil,
			mockRows: sqlmock.NewRows([]string{"id"}),
			wantErr:  ErrNoRows,
		},
		{
			name:    "single row",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow(1, "星期", 24, "三")
				return rows
			}(),
			wantErr: nil,
			wantVal: &TestModel{
				Id:        1,
				FirstName: "星期",
				Age:       24,
				LastName: &sql.NullString{
					String: "三",
					Valid:  true,
				},
			},
		},
		{
			name:    "invalid cols", // 列中存在模型中没有的字段
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "gender"})
				rows.AddRow(1, "星期", 24, "男")
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},
		{
			name:    "many cols", // 列的数量超过模型中定义的字段数量
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "gender"})
				rows.AddRow(1, "星期", "三", 23, "男")
				return rows
			}(),
			wantErr: errs.ErrTooManyColumns,
		},
	}

	// mock
	for _, tt := range tests {
		if tt.mockErr != nil {
			mock.ExpectQuery(tt.query).WillReturnError(tt.mockErr)
		} else {
			mock.ExpectQuery(tt.query).WillReturnRows(tt.mockRows)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})

	}
}

func TestSelector_Select(t *testing.T) {
	db, err := OpenDB(nil)
	require.Nil(t, err)
	tests := []struct {
		name    string
		builder QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "specify columns", // 指定列
			builder: NewSelector[TestModel](db).Select(C("Id"), C("FirstName")),
			want: &Query{
				SQL:  "SELECT `id`,`first_name` FROM `test_models`;",
				Args: nil,
			},
		},
		{
			name:    "specify columns alias", // 指定列别名
			builder: NewSelector[TestModel](db).Select(C("Id").As("_id")).Where(C("Id").Eq(18)),
			want: &Query{
				SQL:  "SELECT `id` AS _id FROM `test_models` WHERE `id` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "aggregate column", // 聚合函数
			builder: NewSelector[TestModel](db).Select(Count("Id")),
			want: &Query{
				SQL:  "SELECT COUNT(`id`) FROM `test_models`;",
				Args: nil,
			},
		},
		{
			name:    "aggregate column alias", // 聚合函数别名
			builder: NewSelector[TestModel](db).Select(Sum("Id").As("sum_id")),
			want: &Query{
				SQL:  "SELECT SUM(`id`) AS sum_id FROM `test_models`;",
				Args: nil,
			},
		},
		{
			name:    "raw sql select", // 原生 SQL Select 查询
			builder: NewSelector[TestModel](db).Select(Raw("DISTINCT `id`")),
			want: &Query{
				SQL:  "SELECT DISTINCT `id` FROM `test_models`;",
				Args: nil,
			},
		},
		{
			name:    "raw sql where",
			builder: NewSelector[TestModel](db).Where(Raw("`name` = ? and `id` = ?", "name", 1).AsPredicate()),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE `name` = ? and `id` = ?;",
				Args: []any{"name", 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkQuerier(b *testing.B) {
	db, err := Open("sqlite3", fmt.Sprintf("file:benchmark_get.db?cache=shared&mode=memory"))
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.db.Exec(TestModel{}.CreateSQL())
	if err != nil {
		b.Fatal(err)
	}

	res, err := db.db.Exec("INSERT INTO `test_models`(`id`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 12, "Wei", 18, "Xin")

	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}
	b.ResetTimer()
	b.Run("unsafe", func(b *testing.B) {
		db.valCreator = valuer.NewUnsafeValuer
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.valCreator = valuer.NewReflectValuer
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
