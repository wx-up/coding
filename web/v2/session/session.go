package session

import (
	"context"
	"net/http"
)

// Session session 中数据的操作
type Session interface {
	Get(ctx context.Context, key string) (string, error)
	// Set value 如果设计为 any 的话，当是 redis 等实现时则需要考虑序列化的问题，比如传递一个 User 接口体
	// 设计成 string 的话，序列化的问题则交给调用方
	Set(ctx context.Context, key string, value string) error
	ID() string
}

// Store 管理 session
type Store interface {
	Generate(ctx context.Context, id string) (Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

// Propagator 管理 sessionID
// 因为是 web 框架，所以接口中直接使用 http.ResponseWriter 和 *http.Request
type Propagator interface {
	Inject(ctx context.Context, id string, resp http.ResponseWriter) error
	Extract(ctx context.Context, req *http.Request) (string, error)
	Remove(ctx context.Context, resp http.ResponseWriter) error
}
