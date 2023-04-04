package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Data{
		Code:    0,
		Message: "成功",
		Data:    data,
	})
}

func Fail(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, Data{
		Code:    -1,
		Message: err.Error(),
	})
}

func AbortFail(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusOK, Data{
		Code:    -1,
		Message: err.Error(),
	})
}

func Abort401(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, Data{
		Code:    -1,
		Message: "请先登录",
	})
}
