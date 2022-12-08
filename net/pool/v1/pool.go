package v1

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Pool interface {
	Get(context.Context) (net.Conn, error)
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

	lock sync.Mutex

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

// Get 获取一个连接
func (d *defaultPool) Get(ctx context.Context) (net.Conn, error) {
	for {
		select {
		// 尝试：从空闲队列中获取一个连接
		case conn := <-d.idleChan:
			// 判读连接是否已经过期
			if conn.lastActive.Add(d.idleTimeout).Before(time.Now()) {
				atomic.AddInt64(&d.cnt, -1)
				_ = conn.c.Close()
				// 结束本次循环，继续下一次循环，尝试获取连接
				continue
			}
			return conn.c, nil
		default:
			// atomic 操作，避免加锁
			atomic.AddInt64(&d.cnt, 1)
			// 当前连接数未超过最大连接数，直接创建一个
			if d.cnt <= d.maxCnt {
				return d.factory()
			}
			// -1
			atomic.AddInt64(&d.cnt, -1)

			// 阻塞获取 Put 的连接
			con := make(chan net.Conn, 1)
			go func() {
				req := &connReq{
					con: make(chan conn, 1),
				}
				d.waitChan <- req
				// 获取可用的连接
				c := <-req.con
				con <- c.c
			}()

			// 超时处理
			select {
			case c := <-con:
				return c, nil
			case <-ctx.Done():
				return nil, errors.New("超时")
			}
		}
	}
}

func (d *defaultPool) Put(c net.Conn) error {
	// 加锁，防止并发
	d.lock.Lock()
	if len(d.waitChan) > 0 {
		connReq := <-d.waitChan
		d.lock.Unlock()
		connReq.con <- conn{
			c:          c,
			lastActive: time.Now(),
		}
		return nil
	}
	d.lock.Unlock()

	// 尝试放入空闲队列，空闲队列满，则 close
	select {
	case d.idleChan <- conn{c: c, lastActive: time.Now()}:
	default:
		atomic.AddInt64(&d.cnt, -1)
		_ = c.Close()
	}
	return nil
}
