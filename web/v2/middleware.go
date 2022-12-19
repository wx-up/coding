package v2

import (
	"fmt"
	"time"
)

type Middleware func(next HandleFunc) HandleFunc

func Cost(next HandleFunc) HandleFunc {
	return func(context *Context) {
		now := time.Now()
		next(context)
		fmt.Printf("cost %d \n", time.Now().Sub(now).Milliseconds())
	}
}

type MiddlewareBuilder struct{}

func (mb *MiddlewareBuilder) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
		}
	}
}
