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
	type fields struct {
		src          any
		dst          any
		ignoreFields []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "src nil",
			fields: fields{
				src: nil,
			},
			wantErr: errors.New("src 不能为 nil"),
		},
		{
			name: "dst nil",
			fields: fields{
				src: &Src{},
				dst: nil,
			},
			wantErr: errors.New("dst 不能为 nil"),
		},
		{
			name: "src type invalid",
			fields: fields{
				src: Src{},
				dst: Dst{},
			},
			wantErr: errors.New("src 只能是指向结构体的指针"),
		},
		{
			name: "dst type invalid",
			fields: fields{
				src: &Src{},
				dst: Dst{},
			},
			wantErr: errors.New("dst 只能是指向结构体的指针"),
		},
		{
			name: "pass",
			fields: fields{
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
			co := &CopyBuilder{
				src:          tt.fields.src,
				dst:          tt.fields.dst,
				ignoreFields: tt.fields.ignoreFields,
			}
			err := co.Builder()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.fields.src.(*Src).Age, tt.fields.dst.(*Dst).Age)
			assert.Equal(t, tt.fields.src.(*Src).Name, tt.fields.dst.(*Dst).Name)
		})
	}
}
