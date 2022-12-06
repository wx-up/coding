package reflect

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetField(t *testing.T) {
	tests := []struct {
		name string

		entity any
		field  string
		newVal any

		wantErr error
	}{
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("参数错误：entity 为 nil"),
		},
		{
			name: "struct",
			entity: struct {
				Name string
			}{
				Name: "wx",
			},
			field:   "Name",
			wantErr: errors.New("参数错误：entity 只能是指向结构体的指针"),
		},
		{
			name:    "pointer not to struct",
			entity:  &(map[string]any{}),
			wantErr: errors.New("参数错误：entity 只能是指向结构体的指针"),
		},
		{
			name: "invalid field",
			entity: &struct {
				Name string
			}{
				Name: "wx",
			},
			field:   "Age",
			newVal:  12,
			wantErr: errors.New("不存在的字段"),
		},
		{
			name: "private field",
			entity: &struct {
				name string
			}{
				name: "wx",
			},
			field:   "name",
			newVal:  "bob",
			wantErr: errors.New("该字段不可修改"),
		},
		{
			name: "pass",
			entity: &User{
				Name: "wx",
			},
			field:   "Name",
			newVal:  "bob",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.entity, tt.field, tt.newVal)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.entity.(*User).Name, tt.newVal)
		})
	}
}

func BenchmarkField(b *testing.B) {
	u := struct {
		Name string
	}{
		Name: "wx",
	}
	refV := reflect.ValueOf(u)

	b.Run("field index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refV.Field(0)
		}
	})

	b.Run("field by name", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refV.FieldByName("Name")
		}
	})
}
