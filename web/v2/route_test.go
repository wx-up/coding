package v2

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_route_add 这里不是要测 add 方法
// 而是要测 add 之后，树的结构是否符合预期
// 所以这里的测试会和之前的有点不一样
func Test_route_add(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			method: http.MethodGet,
			path:   "age",
		},

		{
			method: http.MethodGet,
			path:   "/test/*/info",
		},
	}
	handler := func(ctx *Context) {}

	wantRes := &route{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: handler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: handler,
						children: map[string]*node{
							"detail": {
								path: "detail",
								children: map[string]*node{
									"profile": {
										path:    "profile",
										handler: handler,
									},
								},
							},
						},
					},
					"age": {
						path:    "age",
						handler: handler,
					},
					"test": {
						path: "test",
						starChild: &node{
							path:    "*",
							handler: handler,
						},
					},
				},
			},
		},
	}
	res := &route{
		trees: make(map[string]*node),
	}

	for _, tt := range tests {
		res.add(tt.method, tt.path, handler)
	}

	errStr, ok := wantRes.equal(res)
	assert.True(t, ok, errStr)

	// 函数不能直接比较，所以需要自己写测试方法比较
	// 原理就是使用反射
	// assert.Equal(t, wantRes, res)

	findCases := []struct {
		name   string
		method string
		path   string

		found      bool
		hasHandler bool
		wantPath   string
	}{
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/",

			found:      true,
			wantPath:   "/",
			hasHandler: true,
		},
		{
			name:   "/user",
			method: http.MethodGet,
			path:   "/user",

			found:      true,
			wantPath:   "user",
			hasHandler: true,
		},
		{
			name:   "/user/detail",
			method: http.MethodGet,
			path:   "/user/detail",

			// 注册的是 /user/detail/profile 节点，/user/detail 中间节点应该是找不到的
			found: false,
		},
		{
			name:   "/test/*/info",
			method: http.MethodGet,
			path:   "/test/*/info",

			found:      true,
			wantPath:   "*",
			hasHandler: true,
		},
		{
			name:   "/test/*",
			method: http.MethodGet,
			path:   "/test/*",

			found:      true,
			wantPath:   "*",
			hasHandler: true,
		},
		{
			name:   "/test/ha",
			method: http.MethodGet,
			path:   "/test/ha",

			found:      true,
			wantPath:   "*",
			hasHandler: true,
		},
	}

	for _, fc := range findCases {
		t.Run(fc.name, func(t *testing.T) {
			node, found := res.find(fc.method, fc.path)
			assert.Equal(t, fc.found, found)
			if !found {
				return
			}
			assert.Equal(t, fc.hasHandler, node.handler != nil)
			assert.Equal(t, fc.wantPath, node.path)
		})
	}
}

func (r *route) equal(y *route) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

func Test(t *testing.T) {
	seg := strings.Split("/"[1:], "/")
	fmt.Println(seg)
	fmt.Println(seg[0])
	str := "123"
	fmt.Println(str[3:])
}

func TestFunc(t *testing.T) {
	f := func() {}

	// assert.Equal(t, f, f)

	assert.Equal(t, reflect.ValueOf(f), reflect.ValueOf(f))
}
