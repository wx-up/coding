package session

import (
	"context"
	"net/http"
)

type myPropagator struct {
	cookieName string
	// cookie 存在很多选项，直接提供方法让调用者自己设置，而不是使用字段，它会导致结构体字段很多
	cookieOption func(cookie *http.Cookie)
}

type PropagatorOption func(*myPropagator)

func WithCookieOption(f func(cookie *http.Cookie)) PropagatorOption {
	return func(propagator *myPropagator) {
		propagator.cookieOption = f
	}
}

func WithCookieName(name string) PropagatorOption {
	return func(propagator *myPropagator) {
		propagator.cookieName = name
	}
}
func NewPropagator(opts ...PropagatorOption) Propagator {
	p := &myPropagator{
		cookieName: "sid",
		cookieOption: func(cookie *http.Cookie) {

		},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (m *myPropagator) Inject(ctx context.Context, id string, resp http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  m.cookieName,
		Value: id,
	}
	m.cookieOption(cookie)
	http.SetCookie(resp, cookie)
	return nil
}

func (m *myPropagator) Extract(ctx context.Context, req *http.Request) (string, error) {
	cookie, err := req.Cookie(m.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Remove cookie 的删除其实就是有效期设置为失效
func (m *myPropagator) Remove(ctx context.Context, resp http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   m.cookieName,
		MaxAge: -1,
	}
	m.cookieOption(cookie)
	http.SetCookie(resp, cookie)
	return nil
}
