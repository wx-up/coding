package service

import "net/http"

// serverMux 装饰器模式：新增关闭的时候拒绝连接的功能
type serverMux struct {
	*http.ServeMux
	reject bool
}

func newServerMux() *serverMux {
	return &serverMux{
		ServeMux: http.NewServeMux(),
	}
}

func (sm *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sm.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已经关闭"))
		return
	}
	sm.ServeMux.ServeHTTP(w, r)
}
