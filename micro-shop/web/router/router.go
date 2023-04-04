package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wx-up/coding/micro-shop/web/api"
)

func RegisterRouter(engine *gin.Engine) {
	ug := engine.Group("/user")
	{
		uc := api.NewUserController()
		ug.GET("/list", uc.List)
		ug.POST("/login", uc.Login)
	}
}
