package error

import v2 "github.com/wx-up/coding/web/v2"

type MiddlewareBuilder struct {
	resp map[int][]byte // 错误码和错误信息的映射关系
}

func (mb *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	if mb.resp == nil {
		mb.resp = make(map[int][]byte)
	}
	mb.resp[code] = resp
	return mb
}

func (mb *MiddlewareBuilder) Build() v2.Middleware {
	return func(next v2.HandleFunc) v2.HandleFunc {
		return func(ctx *v2.Context) {
			next(ctx)

			// 如果状态码存在，则将响应替换
			bs, ok := mb.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = bs
			}
		}
	}
}
