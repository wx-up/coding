package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_BindJSON(t *testing.T) {
	type User struct {
		Name string
	}
	type fields struct {
		Req *http.Request
	}
	type args struct {
		val any
	}
	tests := []struct {
		name   string
		fields fields
		args   args

		wantVal any
		wanErr  error
	}{
		{
			name: "happy case",

			fields: fields{
				// 匿名函数的方式创建 request 对象
				Req: func() *http.Request {
					var buffer bytes.Buffer
					buffer.Write([]byte(`{"Name":"Tom"}`))
					req, err := http.NewRequest(http.MethodPost, "/user", &buffer)
					req.Header.Set("Content-Type", "application/json")
					if err != nil {
						t.Fatal(err)
					}
					return req
				}(),
			},

			args: args{
				val: &User{},
			},

			wanErr:  nil,
			wantVal: &User{Name: "Tom"},
		},
		{
			name: "body format not legal",
			fields: fields{
				Req: func() *http.Request {
					var buffer bytes.Buffer
					buffer.Write([]byte(`age=123&name=wx`))
					req, err := http.NewRequest(http.MethodPost, "/user", &buffer)
					if err != nil {
						t.Fatal(err)
					}

					return req
				}(),
			},
			args: args{
				val: &User{},
			},
			wanErr: ErrBodyNotJsonType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Context{
				Req: tt.fields.Req,
			}
			err := c.BindJSON(tt.args.val)
			assert.Equal(t, tt.wanErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantVal, tt.args.val)
		})
	}
}

func TestNumber(t *testing.T) {
	type User struct {
		Age any
	}
	var buffer bytes.Buffer
	buffer.Write([]byte(`{"Age":12}`))
	decoder := json.NewDecoder(&buffer)
	obj := &User{}
	err := decoder.Decode(obj)
	assert.Equal(t, nil, err)
	if err != nil {
		return
	}
	assert.Equal(t, "float64", fmt.Sprintf("%T", obj.Age))

	obj = &User{}
	buffer.Write([]byte(`{"Age":12}`))
	decoder.UseNumber()
	err = decoder.Decode(obj)
	assert.Equal(t, nil, err)
	if err != nil {
		return
	}
	assert.Equal(t, "json.Number", fmt.Sprintf("%T", obj.Age))
}

func TestBuffer(t *testing.T) {
	var b bytes.Buffer
	b.Write([]byte("wx"))

	body := io.NopCloser(&b)
	fmt.Println(io.ReadAll(body))
	fmt.Println(io.ReadAll(body))
	fmt.Println(io.ReadAll(body))
}
