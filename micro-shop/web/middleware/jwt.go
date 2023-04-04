package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wx-up/coding/micro-shop/web/pkg/jwt"
	"github.com/wx-up/coding/micro-shop/web/response"
)

type JwtMiddlewareBuilder struct {
	Secret     string
	ExpireTime int64
}

func (b *JwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			// gin 中间件的方法错误逻辑需要执行 Abort
			// 如果只是 return 的话后面的中间件仍然会执行
			response.Abort401(ctx)
			return
		}

		var opts []jwt.Option
		if b.Secret != "" {
			opts = append(opts, jwt.WithSecret(b.Secret))
		}
		if b.ExpireTime > 0 {
			opts = append(opts, jwt.WithExpireTime(b.ExpireTime))
		}

		res, err := jwt.New(opts...).ParseToken(token)
		if err != nil {
			response.AbortFail(ctx, err)
			return
		}

		ctx.Set("claims", res)

		ctx.Next()
	}
}
