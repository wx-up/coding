package proxy_v2

import "context"

//go:generate mockgen -destination=proxy_gen_test.go -source=proxy.go -package=proxy_v2 Proxy
type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

type Request struct {
	ServiceName string
	MethodName  string

	// 网络传递的时候，尽量都使用字节进行传递
	// 如果这里声明类型为 any，并且赋值为 *GetByIdReq 在服务端的时候会被解析为 map，这时候直接 unmarshal 会报错
	// 因此如果直接声明成 []byte 那么服务端就可以直接 unmarshal
	Arg []byte
}

type Response struct {
	Data []byte
}
