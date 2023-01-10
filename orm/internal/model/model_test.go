package model

import (
	"database/sql"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/wx-up/coding/orm/internal/errs"
	"reflect"
	"strings"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int
	LastName  *sql.NullString
}

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
			val: func() any {
				type TestModel struct {
					Id int64
				}
				return &TestModel{}
			}(),
			want: &Model{
				TableName: "test_models",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
			},
		},
		{
			name: "struct",
			val: func() any {
				type TestModel struct {
					Id int64
				}
				return TestModel{}
			}(),
			want: &Model{
				TableName: "test_models",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
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
			val: func() any {
				type TestModel struct {
					Id int64
				}
				return (*TestModel)(nil)
			}(),
			want: &Model{
				TableName: "test_models",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
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
				TableName: "column_tags",
				FieldMap: map[string]*Field{
					"ID": {
						ColName: "_id",
						Name:    "ID",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
				ColumnMap: map[string]*Field{
					"_id": {
						ColName: "_id",
						Name:    "ID",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
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
				TableName: "column_tags",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Name:    "Id",
						Typ:     reflect.TypeOf(int64(1)),
						Offset:  0,
					},
				},
			},
		},
		{
			name: "custom table name", // 自定义表名
			val:  &CustomTableName{},
			want: &Model{
				TableName: "users",
				FieldMap:  map[string]*Field{},
				ColumnMap: map[string]*Field{},
			},
		},
		{
			name: "empty table name",
			val:  &EmptyTableName{}, // TableName 返回空字符串的时候，则走默认规则
			want: &Model{
				TableName: "empty_table_names",
				FieldMap:  map[string]*Field{},
				ColumnMap: map[string]*Field{},
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
			assert.Equal(t, tt.want.ColumnMap, got.ColumnMap)
			assert.Equal(t, tt.want.FieldMap, got.FieldMap)
			assert.Equal(t, tt.want.TableName, got.TableName)
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
