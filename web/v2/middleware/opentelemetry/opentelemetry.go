package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"

	v2 "github.com/wx-up/coding/web/v2"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() v2.Middleware {
	return func(next v2.HandleFunc) v2.HandleFunc {
		return func(ctx *v2.Context) {
			// 使用 request 的 context 作为父context
			// 后续请求中的 span 都基于这个 context 生成
			spanCtx, span := m.Tracer.Start(ctx.Req.Context(), "web")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))

			// 只取 path 的一段，如果全部存储的话，当有人攻击时会给 trace 服务造成很大的压力
			// 比如 path 输入很长
			path := ctx.Req.URL.Path
			if len(path) >= 256 {
				path = path[:256]
			}
			span.SetAttributes(attribute.String("http.path", path))

			// 将 request 中的 context 设置为 spanCtx
			// 注意：这里复制 request 是一个性能很差的事情
			// 优化：可以在 ctx 中新增一个 Ctx 字段
			ctx.Req = ctx.Req.WithContext(spanCtx)

			next(ctx)

			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
