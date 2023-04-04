package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorsMiddlewareBuilder struct{}

func (m *CorsMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		method := ctx.Request.Method

		ctx.Header("Access-Control-Allow-Origin", origin)
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token,Authorization,x-token")
		ctx.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		}
	}
}
