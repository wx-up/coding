package reflect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
}

func Test_iterateVal(t *testing.T) {
	u1 := &User{Name: "星期三"}
	u2 := &u1
	tests := []struct {
		// 名称
		name string

		// 输入部分
		val interface{}

		// 输出部分
		wantRes map[string]any
		wantErr error
	}{
		{
			name:    "nil",
			val:     nil,
			wantErr: errors.New("val 不能为 nil"),
		},
		{
			name:    "user",
			val:     User{Name: "wx"},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "wx",
			},
		},
		{
			name:    "pointer",
			val:     &User{Name: "bob"},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "bob",
			},
		},
		{
			name:    "multiple pointer",
			val:     u2,
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "星期三",
			},
		},
		// 非法输入
		{
			name:    "slice",
			val:     []int64{1, 2},
			wantErr: errors.New("非法输入"),
		},
		// 非法指针输入
		{
			name:    "pointer to map",
			val:     &(map[string]interface{}{}),
			wantErr: errors.New("非法输入"),
		},

		// 结构体字段是结构体或者指针直接忽略
		{
			name: "filed is struct",
			val: struct {
				Name  string
				Hobby struct {
					Hobby string
				}
			}{
				Name: "星期三",
				Hobby: struct {
					Hobby string
				}{
					Hobby: "篮球",
				},
			},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "星期三",
			},
		},

		// 结构体包含非导出字段
		{
			name: "filed is unexported",
			val: struct {
				Name string
				age  int
			}{
				Name: "星期三",
				age:  27,
			},
			wantRes: map[string]any{
				"Name": "星期三",
				"age":  0,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iterateVal(tt.val)
			assert.Equal(t, tt.wantErr, err)

			// 如果出现错误，后续的值断言就没有意义了
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantRes, got)
		})
	}
}
