package v2

import (
	"net/http"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	server.Get("/cookie", func(ctx *Context) {
		cookie := &http.Cookie{
			Name:     "test_cookie",
			Value:    "value",
			HttpOnly: true,
		}
		http.SetCookie(ctx.Resp, cookie)
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("你好呀")
		return
	})
	server.Start(":8080")
}
