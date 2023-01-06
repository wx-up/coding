package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int
	LastName  *sql.NullString
}

func TestSelect_Build(t *testing.T) {
	db, err := NewDB()
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
