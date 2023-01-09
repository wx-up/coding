package model

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/wx-up/coding/orm/internal/errs"
	"strings"
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
		{
			name: "column tag",
			val: func() any { // 当值比较复杂的时候可以使用函数
				type ColumnTag struct {
					ID int64 `orm:"column=_id"`
				}
				return &ColumnTag{}
			}(),
			wantErr: nil,
			want: &Model{
				tableName: "column_tags",
				fieldMap: map[string]*field{
					"ID": {
						colName: "_id",
					},
				},
			},
		},
		{
			name: "invalid column tag",
			val: func() any {
				type ColumnTag struct {
					Id int64 `orm:"column="` // 按照默认的解析规则
				}
				return &ColumnTag{}
			}(),
			want: &Model{
				tableName: "column_tags",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
				},
			},
		},
		{
			name: "custom table name", // 自定义表名
			val:  &CustomTableName{},
			want: &Model{
				tableName: "users",
				fieldMap:  map[string]*field{},
			},
		},
		{
			name: "empty table name",
			val:  &EmptyTableName{}, // TableName 返回空字符串的时候，则走默认规则
			want: &Model{
				tableName: "empty_table_names",
				fieldMap:  map[string]*field{},
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

func Test(t *testing.T) {
	fmt.Println(strcase.ToSnake("Id"))

	s := "na"
	fmt.Println(strings.SplitN(s, "=", 2))
}

type CustomTableName struct {
}

func (*CustomTableName) TableName() string {
	return "users"
}

type EmptyTableName struct {
}

func (*EmptyTableName) TableName() string {
	return ""
}
