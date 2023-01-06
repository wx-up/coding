package errs

import (
	"errors"
	"fmt"
)

var (
	ErrParseModelValType = errors.New("orm: 只支持结构体或者结构体的指针")
)

func NewErrUnsupportedExpressionType(expr any) error {
	return fmt.Errorf("orm: 不支持的表达式 %v", expr)
}
