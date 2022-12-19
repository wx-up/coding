package v2

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server struct {
	route

	ms []Middleware
}

func NewServer() *Server {
	return &Server{
		route: newRoute(),
	}
}

func (s *Server) Use(ms ...Middleware) {
	s.ms = append(s.ms, ms...)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	// 将路由的查找和业务的执行，包装成责任链中的最后一环
	root := s.serve

	// 构建洋葱模型，从尾部往头部
	for i := len(s.ms) - 1; i >= 0; i-- {
		root = s.ms[i](root)
	}

	// 将响应刷到前端的中间件放到第一个
	var flushMiddleware Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)

			s.flush(ctx)
		}
	}
	root = flushMiddleware(root)
	// 执行
	root(ctx)
}

func (s *Server) flush(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	_, _ = ctx.Resp.Write(ctx.RespData)
}

func (s *Server) serve(ctx *Context) {
	node, ok := s.find(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("路由未找到")
		return
	}

	node.handler(ctx)
}

func (s *Server) Get(path string, handler HandleFunc) {
	s.AddRoute(http.MethodGet, path, handler)
}

func (s *Server) AddRoute(method string, path string, handler HandleFunc) {
	s.add(method, path, handler)
}
