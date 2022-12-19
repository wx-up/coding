package error

import (
	"net/http"
	"testing"

	v2 "github.com/wx-up/coding/web/v2"
)

func Test(t *testing.T) {
	serve := v2.NewServer()

	serve.Get("/get", func(ctx *v2.Context) {
		ctx.RespData = []byte("hello world")
		ctx.RespStatusCode = http.StatusOK
	})

	serve.Use((&MiddlewareBuilder{}).RegisterError(http.StatusNotFound, []byte("路由不存在，哈哈哈")).Build())

	serve.Start(":8081")
}
