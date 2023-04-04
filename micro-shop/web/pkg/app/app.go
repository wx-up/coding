package app

import (
	"github.com/gin-gonic/gin"
	"github.com/wx-up/coding/micro-shop/web/pkg/jwt"
)

// CurrentUid 获取当前登录的用户id
func CurrentUid(ctx *gin.Context) uint64 {
	res, ok := ctx.Get("claims")
	if !ok {
		return 0
	}
	if claims, ok := res.(jwt.CustomInfo); ok {
		return claims.Id
	}
	return 0
}
