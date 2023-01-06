package orm

import (
	"github.com/stretchr/testify/assert"
	"github.com/wx-up/coding/orm/internal/errs"
	"testing"
)

func Test_parseModel(t *testing.T) {
	/*
		type TestModel struct {
			Id        int64
			FirstName string
			Age       int
			LastName  *sql.NullString
		}
	*/
	tests := []struct {
		name    string
		val     any
		want    *Model
		wantErr error
	}{
		{
			name: "ptr",
			val:  &TestModel{},
			want: &Model{
				tableName: "test_models",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name: "struct",
			val:  TestModel{},
			want: &Model{
				tableName: "test_models",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errs.ErrParseModelValType,
		},
		{
			name:    "nil",
			val:     nil,
			wantErr: errs.ErrParseModelValType,
		},
		{
			name: "nil with type",
			// 这种情况是可以获取到信息的，区别于 nil
			val: (*TestModel)(nil),
			want: &Model{
				tableName: "test_models",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Registry{}
			got, err := r.parseModel(tt.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
