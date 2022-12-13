package v1

import (
	"net"
	"net/http"
)

type HandleFunc func(*Context)

type HTTPServer struct{}

// Start 使用 http.Serve 来启动，获取更大的灵活性，如将端口监听和服务器启动分离等等
func (m *HTTPServer) Start(addr string) error {
	// 端口启动前
	// 处理一些钩子函数
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	// 端口启动后
	// 处理一些钩子函数，比如将本服务注册到管理平台，比如 etcd，然后打开管理界面就可以看到这个服务

	return http.Serve(listener, m)
}

func (m *HTTPServer) Get(path string, handler HandleFunc) {
	m.AddRoute(http.MethodGet, path, handler)
}

func (m *HTTPServer) AddRoute(method string, path string, handler HandleFunc) {
}

func (m *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	m.serve(ctx)
}

func (m *HTTPServer) serve(ctx *Context) {
}

// Start1 该实现方式相当于将 start 的两步合并了，做不了一些钩子的功能
func (m *HTTPServer) Start1(addr string) error {
	return http.ListenAndServe(":8081", m)
}

func main() {
	// AddRoute 只支持单个 HandleFunc，如果需要实现多个 HandleFunc 则自己组合
	var h1 HandleFunc
	var h2 HandleFunc
	server := &HTTPServer{}
	server.AddRoute(http.MethodGet, "/user", func(context *Context) {
		h1(context)
		h2(context)
	})
}
