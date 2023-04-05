package proxy_v2

import (
	"context"

	"github.com/wx-up/coding/micro/custom_protocol/message"
)

//go:generate mockgen -destination=proxy_gen_test.go -source=proxy.go -package=proxy_v2 Proxy
type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}
