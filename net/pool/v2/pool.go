package v2

import (
	"context"
	"net"
	"sync"
	"time"
)

type Factory func() (net.Conn, error)

type Pool struct {
	// 空闲连接队列
	idleConnQueue chan *idleConn

	// 请求连接队列
	reqConnQueue []*connReq

	// 最大连接数
	maxCnt int
	// 当前连接数
	cnt int

	// 最大空闲时间
	maxIdleTime time.Duration

	// 构建连接
	factory Factory

	// 初始化容量
	initCap int

	lock sync.Mutex
}

type Option func(*Pool)

func WithMaxCnt(cnt int) Option {
	return func(pool *Pool) {
		pool.maxCnt = cnt
	}
}

func WithMaxIdleTime(t time.Duration) Option {
	return func(pool *Pool) {
		pool.maxIdleTime = t
	}
}

func WithMaxIdle(idle int) Option {
	return func(pool *Pool) {
		pool.idleConnQueue = make(chan *idleConn, idle)
	}
}

func WithInitCap(cap int) Option {
	return func(pool *Pool) {
		pool.initCap = cap
	}
}

func New(factory Factory, opts ...Option) *Pool {
	pool := &Pool{
		idleConnQueue: make(chan *idleConn, 80),
		maxCnt:        128,
		maxIdleTime:   time.Minute * 30,
		factory:       factory,
	}
	for _, opt := range opts {
		opt(pool)
	}

	// 初始容量大于最大空闲连接数量
	if pool.initCap > cap(pool.idleConnQueue) {
		panic("micro：初始容量大于最大空闲连接数量")
	}

	// 将时间放到外面初始化，相对于 lastActiveTime 来说是有误差的
	// 毕竟连接的创建需要三次握手
	timeNow := time.Now()
	for i := 0; i < pool.initCap; i++ {
		conn, err := pool.factory()
		if err != nil {
			panic(err)
		}
		pool.idleConnQueue <- &idleConn{
			conn:           conn,
			lastActiveTime: timeNow,
		}
	}

	return pool
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	// 超时判断
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default: // 需要有一个 default 分支，否则会一直阻塞
	}

	for {
		select {
		case conn := <-p.idleConnQueue: // 空闲队列存在可用连接
			// 本次拿到的连接如果已经过期了，则继续尝试获取
			if conn.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = conn.conn.Close()
				p.lock.Lock()
				p.cnt--
				p.lock.Unlock()
				continue
			}

			return conn.conn, nil
		default: // 空闲队列不存在可用连接
			p.lock.Lock()
			if p.cnt >= p.maxCnt { // 超过最大连接数，需要阻塞等待别人归还
				reqConn := &connReq{connChan: make(chan net.Conn, 1)}
				p.reqConnQueue = append(p.reqConnQueue, reqConn)
				// 解锁，因此对 reqConnQueue 的修改已经完成
				p.lock.Unlock()
				select {
				case <-ctx.Done(): // 超时了
					// 选项一：从 reqConnQueue 将自己删除掉 + 转发
					// 选项二：直接使用转发（ Put 的实现本身就带有删除 reqConnQueue 的逻辑 ）
					// 这里使用转发
					go func() {
						conn := <-reqConn.connChan
						_ = p.Put(context.Background(), conn)
					}()

					return nil, ctx.Err()
				case conn := <-reqConn.connChan: // 别人归还了连接，直接返回，最好可以 ping 一下
					return conn, nil
				}
			}

			// 没有超过最大连接数
			conn, err := p.factory()
			if err != nil {
				return nil, err
			}
			p.cnt++
			p.lock.Unlock()
			return conn, nil
		}
	}
}

func (p *Pool) Put(ctx context.Context, conn net.Conn) error {
	p.lock.Lock()

	// 存在阻塞的请求
	if len(p.reqConnQueue) > 0 {
		// 这里从队首开始拿
		// 正常应该考虑从队尾开始拿，这样子拿到的 reqConn 在 Get 中超时的概率比较少
		reqConn := p.reqConnQueue[0]
		p.reqConnQueue = p.reqConnQueue[1:]
		p.lock.Unlock()
		reqConn.connChan <- conn
		return nil
	}

	p.lock.Unlock()

	// 没有阻塞的请求
	select {
	case p.idleConnQueue <- &idleConn{conn: conn, lastActiveTime: time.Now()}: // 尝试放入空闲队列
	default: // 空闲队列满了
		_ = conn.Close()
		p.lock.Lock()
		p.cnt--
		p.lock.Unlock()
	}
	return nil
}

type idleConn struct {
	conn net.Conn
	// 上一次使用的时间
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}
