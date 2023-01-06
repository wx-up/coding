package errs

import (
	"errors"
	"fmt"
)

var (
	ErrParseModelValType = errors.New("orm: 只支持结构体或者结构体的指针")
)

func NewErrUnsupportedExpressionType(expr any) error {
	// 错误信息加上 orm 前缀，标识错误的源头
	return fmt.Errorf("orm: 不支持的表达式 %v", expr)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm: 未知字段 %s", name)
}
