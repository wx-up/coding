package access_log

import (
	"encoding/json"
	"log"

	v2 "github.com/wx-up/coding/web/v2"
)

type MiddlewareBuilder struct {
	logFunc func(content string) error
}

// SetLogFunc 可以重置输出
func (mb *MiddlewareBuilder) SetLogFunc(f func(string) error) *MiddlewareBuilder {
	mb.logFunc = f
	return mb
}

func New() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		// 默认在控制台输出
		logFunc: func(content string) error {
			log.Println(content)
			return nil
		},
	}
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

func (mb *MiddlewareBuilder) Builder() v2.Middleware {
	return func(next v2.HandleFunc) v2.HandleFunc {
		return func(ctx *v2.Context) {
			defer func() {
				if mb.logFunc == nil {
					return
				}
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchPath,
					Path:       ctx.Req.URL.Path,
					HTTPMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				_ = mb.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}
