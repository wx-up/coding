package service

import (
	"context"
	"net/http"
)

type Server struct {
	name string
	mux  *serverMux
	// srv 实例用于启动服务和关闭服务
	srv *http.Server
}

func NewServer(name string, addr string) *Server {
	mux := newServerMux()
	return &Server{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Handle 注册路由
func (s *Server) Handle(pattern string, handler http.HandlerFunc) {
	s.mux.Handle(pattern, handler)
}

// rejectReq 拒绝请求
func (s *Server) rejectReq() {
	s.mux.reject.Store(true)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.Background())
}

func (s *Server) StopWithContext(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
