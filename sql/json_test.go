package sql

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonColumn_Value(t *testing.T) {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"Name":"Tom"}`), value)
	js = JsonColumn[User]{}
	value, err = js.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		wantErr error
		wantVal User
	}{
		{
			name:    "nil",
			wantErr: errors.New("src 不能为 nil"),
		},
		{
			name:    "string",
			src:     `{"Name":"Tom"}`,
			wantVal: User{Name: "Tom"},
		},
		{
			name:    "bytes",
			src:     []byte(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
			assert.True(t, js.Valid)
		})
	}
}

func TestJsonColumn_ScanTypes(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		jsonSlice := &JsonColumn[[]string]{}
		err := jsonSlice.Scan(`["a","b","c"]`)
		assert.Nil(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, jsonSlice.Val)
		assert.True(t, jsonSlice.Valid)
		if err != nil {
			return
		}
		res, err := jsonSlice.Value()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`["a","b","c"]`), res)
	})

	t.Run("map", func(t *testing.T) {
		jsonMap := &JsonColumn[map[string]string]{}
		err := jsonMap.Scan(`{"name":"wx","age":"12"}`)
		assert.Nil(t, err)
		assert.Equal(t, map[string]string{"name": "wx", "age": "12"}, jsonMap.Val)
		assert.True(t, jsonMap.Valid)
		if err != nil {
			return
		}
		res, err := jsonMap.Value()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"age":"12","name":"wx"}`), res)
	})
}

type User struct {
	Name string
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(value.([]byte)))
	// Output:
	// {"Name":"Tom"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan(`{"Name":"Tom"}`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(js.Val)
	// Output:
	// {Tom}
}
