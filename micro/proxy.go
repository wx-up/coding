package micro

import "context"

type Request struct {
	ServerName string
	MethodName string
	Arg        any // 暂时不考虑 context 参数
}

type Response struct {
	Data []byte
}

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}
