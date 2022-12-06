package reflect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Src struct {
	Name string
	Age  int64
}
type Dst struct {
	Age  int64
	Name string
}

func TestCopyBuilder_Builder(t *testing.T) {
	type args struct {
		src          any
		dst          any
		ignoreFields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "src nil",
			args: args{
				src: nil,
			},
			wantErr: errors.New("src 不能为 nil"),
		},
		{
			name: "dst nil",
			args: args{
				src: &Src{},
				dst: nil,
			},
			wantErr: errors.New("dst 不能为 nil"),
		},
		{
			name: "src type invalid",
			args: args{
				src: Src{},
				dst: Dst{},
			},
			wantErr: errors.New("src 只能是指向结构体的指针"),
		},
		{
			name: "dst type invalid",
			args: args{
				src: &Src{},
				dst: Dst{},
			},
			wantErr: errors.New("dst 只能是指向结构体的指针"),
		},
		{
			name: "pass",
			args: args{
				src: &Src{
					Name: "wx",
					Age:  16,
				},
				dst: &Dst{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.args.src, tt.args.dst, tt.args.ignoreFields)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.args.src.(*Src).Age, tt.args.dst.(*Dst).Age)
			assert.Equal(t, tt.args.src.(*Src).Name, tt.args.dst.(*Dst).Name)
		})
	}
}
