package request

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Validator interface {
	Validate() map[string][]string
}

func Validate(ctx *gin.Context, in Validator) error {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return errors.New("第二个参数必须是指针")
	}

	if err := ctx.ShouldBind(in); err != nil {
		return err
	}

	res := in.Validate()
	return getFirstError(res)
}

func getFirstError(res map[string][]string) error {
	if len(res) <= 0 {
		return nil
	}
	for k, v := range res {
		if len(v) <= 0 {
			continue
		}
		return fmt.Errorf("参数错误：%s - %s", k, v[0])
	}
	return nil
}
