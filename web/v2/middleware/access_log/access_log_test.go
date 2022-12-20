package access_log

import (
	"net/http"
	"testing"

	v2 "github.com/wx-up/coding/web/v2"
)

func Test(t *testing.T) {
	srv := v2.NewServer()
	srv.Get("/user", func(ctx *v2.Context) {
		ctx.RespData = []byte("hello world")
		ctx.RespStatusCode = http.StatusOK
	})

	builder := New()
	srv.Use(builder.Builder())
	srv.Start(":8081")
}
