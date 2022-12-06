package reflect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_iterate(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name string

		// 输入
		args args

		// 输出
		want    []any
		wantErr error
	}{
		{
			name: "nil",
			args: args{
				val: nil,
			},
			wantErr: errors.New("val 不能为 nil"),
		},
		{
			name: "slice",
			args: args{
				val: []int{1, 2, 3},
			},
			wantErr: nil,
			want:    []any{1, 2, 3},
		},
		{
			name: "array",
			args: args{
				val: [3]int{5, 6, 7},
			},
			wantErr: nil,
			want:    []any{5, 6, 7},
		},
		{
			name: "string",
			args: args{
				val: "abc",
			},
			want:    []any{uint8('a'), uint8('b'), uint8('c')},
			wantErr: nil,
		},
		{
			name: "invalid val",
			args: args{
				val: 18,
			},
			wantErr: errors.New("val 类型错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iterate(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_iterateMap(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name    string
		args    args
		wantKey []any
		wantVal []any
		wantErr error
	}{
		{
			name: "nil",
			args: args{
				val: nil,
			},
			wantErr: nil,
		},
		{
			name: "invalid val",
			args: args{
				val: "abc",
			},
			wantErr: errors.New("val 类型错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := iterateMap(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			if !assert.Equal(t, tt.wantKey, got1) {
				return
			}
			assert.Equal(t, tt.wantVal, got2)
		})
	}
}

func Test_iterateMap1(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name       string
		args       args
		wantKeys   []any
		wantValues []any
		wantErr    error
	}{
		{
			name: "val valid",
			args: args{
				val: &map[string]string{
					"name": "bob",
				},
			},
			wantErr: errors.New("val 类型错误"),
		},
		{
			name: "pass",
			args: args{
				val: map[string]string{
					"name":  "wx",
					"hobby": "basketball",
				},
			},
			wantKeys:   []any{"name", "hobby"},
			wantValues: []any{"wx", "basketball"},
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := iterateMap(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantKeys, got)
			assert.Equal(t, tt.wantValues, got1)
		})
	}
}
