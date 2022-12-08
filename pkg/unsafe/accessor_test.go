package unsafe

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	fmt.Println(strings.SplitN("date ", " ", 2))
}

func TestAccessor_Field(t *testing.T) {
	type args struct {
		name   string
		entity interface{}
	}
	tests := []struct {
		name string

		// 输入参数
		args args

		want    any
		wantErr error
	}{
		{
			name: "field name empty",
			args: args{
				name: "",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: errors.New("name 不能为空"),
		},
		{
			name: "field not exist",
			args: args{
				name: "Age",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: errors.New("不存在的字段"),
		},
		{
			name: "pass",
			args: args{
				name: "Name",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: nil,
			want:    "wx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			access, err := NewAccessor(tt.args.entity)
			assert.Equal(t, nil, err)
			if err != nil {
				return
			}
			got, err := access.Field(tt.args.name)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccessor_SetField(t *testing.T) {
	type args struct {
		name   string
		val    any
		entity any
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "name empty",
			args: args{
				name: "",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: errors.New("name 不能为空"),
		},
		{
			name: "val nil",
			args: args{
				name: "Name",
				val:  nil,
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: errors.New("val 不能为nil"),
		},
		{
			name: "field not exit",
			args: args{
				name: "Age",
				val:  "bob",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: errors.New("不存在的字段"),
		},
		{
			name: "pass",
			args: args{
				name: "Name",
				val:  "Bob",
				entity: &struct {
					Name string
				}{
					Name: "wx",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			access, err := NewAccessor(tt.args.entity)
			assert.Equal(t, nil, err)
			if err != nil {
				return
			}
			err = access.SetField(tt.args.name, tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.args.val, tt.args.entity.(*struct{ Name string }).Name)
		})
	}
}
