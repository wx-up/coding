package v1

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := New(func() (net.Conn, error) {
		return &mockConn{}, nil
	}, WithMaxCnt(3), WithMaxIdle(2))

	c1, err := pool.Get(context.Background())
	assert.Nil(t, err)
	c2, err := pool.Get(context.Background())
	assert.Nil(t, err)
	c3, err := pool.Get(context.Background())
	assert.Nil(t, err)

	// 正常放回去
	err = pool.Put(c1)
	assert.Nil(t, err)
	err = pool.Put(c2)
	assert.Nil(t, err)

	// 空闲队列满了，c3 会关闭
	err = pool.Put(c3)
	assert.Nil(t, err)
	assert.True(t, c3.(*mockConn).closed)
}

func TestPool_BlockTimeout(t *testing.T) {
	pool := New(func() (net.Conn, error) {
		return &mockConn{}, nil
	}, WithMaxCnt(3), WithMaxIdle(2))

	_, err := pool.Get(context.Background())
	assert.Nil(t, err)
	_, err = pool.Get(context.Background())
	assert.Nil(t, err)
	_, err = pool.Get(context.Background())
	assert.Nil(t, err)

	// 阻塞
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// 超时
	_, err = pool.Get(timeout)
	assert.Equal(t, errors.New("超时"), err)
}

func TestPool_GetBlock(t *testing.T) {
	pool := New(func() (net.Conn, error) {
		return &mockConn{}, nil
	}, WithMaxCnt(3), WithMaxIdle(2))

	c1, err := pool.Get(context.Background())
	assert.Nil(t, err)
	_, err = pool.Get(context.Background())
	assert.Nil(t, err)
	_, err = pool.Get(context.Background())
	assert.Nil(t, err)

	// 两秒后将c1放回去
	go func() {
		time.Sleep(time.Second * 2)
		_ = pool.Put(c1)
	}()

	// 这里拿到的应该是 c1
	c2, err := pool.Get(context.Background())
	assert.Nil(t, err)
	// 断言相等
	assert.Equal(t, c1, c2)
}

type mockConn struct {
	closed bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) RemoteAddr() net.Addr {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) SetDeadline(t time.Time) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	// TODO implement me
	panic("implement me")
}
