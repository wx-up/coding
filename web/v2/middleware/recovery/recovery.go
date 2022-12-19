package recovery

import v2 "github.com/wx-up/coding/web/v2"

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *v2.Context)
}

func (mb *MiddlewareBuilder) Build() v2.Middleware {
	return func(next v2.HandleFunc) v2.HandleFunc {
		return func(ctx *v2.Context) {
			defer func() {
				if err := recover(); err != nil {
					// 当 panic 之后，返回指定的响应
					ctx.RespData = []byte(mb.ErrMsg)
					ctx.RespStatusCode = mb.StatusCode

					// 万一 LogFunc 也 panic 那就无能为力
					if mb.LogFunc != nil {
						mb.LogFunc(ctx)
					}
				}
			}()

			next(ctx)
		}
	}
}
