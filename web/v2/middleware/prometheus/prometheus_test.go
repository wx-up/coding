package prometheus

import (
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	v2 "github.com/wx-up/coding/web/v2"
)

func Test(t *testing.T) {
	s := v2.NewServer()
	s.Get("/", func(ctx *v2.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *v2.Context) {
		time.Sleep(time.Second)
		ctx.RespData = []byte(`hello world`)
		ctx.RespStatusCode = http.StatusOK
	})

	s.Use((&MiddlewareBuilder{
		Subsystem: "web",
		Name:      "http_request",
		Help:      "这是测试例子",
		ConstLabels: map[string]string{
			"instance_id": "1234567",
		},
	}).Build())
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 一般来说，在实际中我们都会单独准备一个端口给监控
		http.ListenAndServe(":2112", nil)
	}()
	s.Start(":8081")
}
