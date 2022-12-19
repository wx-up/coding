package opentelemetry

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel"

	v2 "github.com/wx-up/coding/web/v2"
)

func TestMiddleware(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("")
	initJaeger(t)
	s := v2.NewServer()
	s.Get("/", func(ctx *v2.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *v2.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		c, third2 := tracer.Start(c, "third_layer_1")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()
		ctx.RespStatusCode = 200
		ctx.RespData = []byte("hello, world")
	})

	s.Use((&MiddlewareBuilder{Tracer: tracer}).Build())
	s.Start(":8081")
}
