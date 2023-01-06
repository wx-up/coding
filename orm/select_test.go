package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
}

func TestSelect_Build(t *testing.T) {
	tests := []struct {
		name    string
		builder QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "from",
			builder: NewSelector[TestModel]().From("test_model_tbl"),
			want: &Query{
				SQL: "SELECT * FROM test_model_tbl;",
			},
			wantErr: nil,
		},
		{
			name:    "from with quotation mark", // 反引号
			builder: NewSelector[TestModel]().From("`test_model_tbl`"),
			want: &Query{
				SQL: "SELECT * FROM `test_model_tbl`;",
			},
			wantErr: nil,
		},
		{
			name:    "no from",
			builder: NewSelector[TestModel](),
			want: &Query{
				SQL: "SELECT * FROM `test_models`;",
			},
			wantErr: nil,
		},
		{
			name:    "from but empty",
			builder: NewSelector[TestModel]().From(""),
			want: &Query{
				SQL: "SELECT * FROM `test_models`;",
			},
			wantErr: nil,
		},
		{
			name:    "from with db",
			builder: NewSelector[TestModel]().From("byn.test_model"),
			want: &Query{
				SQL: "SELECT * FROM byn.test_model;",
			},
			wantErr: nil,
		},
		{
			name:    "single predicate",
			builder: NewSelector[TestModel]().Where(C("id").Eq(12)),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE `id` = ?;",
				Args: []any{12},
			},
			wantErr: nil,
		},
		{
			name:    "and predicates",
			builder: NewSelector[TestModel]().Where(C("age").Gt(12).And(C("age").Lt(24))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{12, 24},
			},
			wantErr: nil,
		},
		{
			name:    "or predicates",
			builder: NewSelector[TestModel]().Where(C("name").Eq("bob").Or(C("age").Gt(12))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`name` = ?) OR (`age` > ?);",
				Args: []any{"bob", 12},
			},
		},
		{
			name:    "not predicates",
			builder: NewSelector[TestModel]().Where(C("name").Eq("bob").And(Not(C("age").Eq(12)))),
			want: &Query{
				SQL:  "SELECT * FROM `test_models` WHERE (`name` = ?) AND (NOT (`age` = ?));",
				Args: []any{"bob", 12},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Build()
			assert.Nil(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
