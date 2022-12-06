package reflect

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Order struct {
	buy   int64
	count int64
}

func (o Order) Buyer() int64 {
	return o.buy
}

func (o *Order) Count() int64 {
	return o.count
}

func Test_iterateFunc(t *testing.T) {
	tests := []struct {
		name string

		val any

		want    map[string]*FuncInfo
		wantErr error
	}{
		{
			name:    "nil",
			val:     nil,
			wantErr: errors.New("val 不能为 nil"),
		},
		{
			name:    "invalid val",
			val:     map[string]string{},
			wantErr: errors.New("不支持的 val 类型"),
		},
		{
			name: "struct",
			val: Order{
				buy: 64,
			},
			wantErr: nil,
			want: map[string]*FuncInfo{
				"Buyer": {
					Name:   "Buyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(64)},
				},
			},
		},
		{
			name: "struct pointer",
			val: &Order{
				buy:   64,
				count: 100,
			},
			wantErr: nil,
			want: map[string]*FuncInfo{
				"Buyer": {
					Name:   "Buyer",
					In:     []reflect.Type{reflect.TypeOf(&Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(64)},
				},
				"Count": {
					Name:   "Count",
					In:     []reflect.Type{reflect.TypeOf(&Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(100)},
				},
			},
		},
		{
			name: "struct and pointer method",
			val: Order{
				count: 64,
			},
			wantErr: nil,
			want: map[string]*FuncInfo{
				"Buyer": { // 结构体不存在指针方法
					Name:   "Buyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(0)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iterateFunc(tt.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
