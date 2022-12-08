package v1

import (
	"net"
	"sync"
	"time"
)

type Pool interface {
	Get() (net.Conn, error)
	Put(net.Conn) error
}

type Factory func() (net.Conn, error)

type conn struct {
	c          net.Conn
	lastActive time.Time
}

type connReq struct {
	con chan conn
}

type defaultPool struct {
	cnt    int64 // 当前连接数
	maxCnt int64 // 最大连接数

	waitChan chan *connReq // 阻塞队列
	idleChan chan conn     // 空闲连接

	mutex sync.Mutex

	factory Factory // 创建连接

	idleTimeout time.Duration
}

type Option func(*defaultPool)

func WithMaxIdle(maxIdle int64) Option {
	return func(p *defaultPool) {
		p.idleChan = make(chan conn, maxIdle)
	}
}

func WithMaxCnt(maxCnt int64) Option {
	return func(p *defaultPool) {
		p.maxCnt = maxCnt
	}
}

var _ Pool = (*defaultPool)(nil)

func New(f Factory, opts ...Option) Pool {
	p := &defaultPool{
		maxCnt:      128,
		waitChan:    make(chan *connReq, 128),
		idleChan:    make(chan conn, 16),
		factory:     f,
		idleTimeout: time.Second * 30,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (d *defaultPool) Get() (net.Conn, error) {
	panic("implement me")
}

func (d *defaultPool) Put(n net.Conn) error {
	panic("implement me")
}
