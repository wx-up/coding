package recovery

import (
	"net/http"
	"testing"

	v2 "github.com/wx-up/coding/web/v2"
)

func Test(t *testing.T) {
	serve := v2.NewServer()

	serve.Get("/get", func(ctx *v2.Context) {
		panic("我 panic 了")
	})

	serve.Use((&MiddlewareBuilder{
		StatusCode: http.StatusInternalServerError,
		ErrMsg:     "出现 panic 了",
	}).Build())

	serve.Start(":8081")
}
