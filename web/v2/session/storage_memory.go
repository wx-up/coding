package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// myStore 管理 session 实例
type myStore struct {
	// 利用一个内存缓存来帮助我们管理过期时间
	c *cache.Cache
	// 缓存过期时间
	expiration time.Duration
	// 每隔多久清理过期的缓存
	cleanupInterval time.Duration
}

type StoreOption func(*myStore)

func WithStoreCleanupInterval(d time.Duration) StoreOption {
	return func(store *myStore) {
		store.cleanupInterval = d
	}
}

func NewStore(expire time.Duration, opts ...StoreOption) Store {
	s := &myStore{
		expiration:      expire,
		cleanupInterval: time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.c = cache.New(s.expiration, s.cleanupInterval)
	return s
}

func (m *myStore) Generate(ctx context.Context, id string) (Session, error) {
	session := &memorySession{
		data: make(map[string]string),
		id:   id,
	}

	// 放入内存缓存
	m.c.Set(id, session, 0)

	return session, nil
}

// Refresh 刷新有效期
func (m *myStore) Refresh(ctx context.Context, id string) error {
	session, err := m.Get(ctx, id)
	if err != nil {
		return err
	}
	m.c.Set(session.ID(), session, m.expiration)

	return nil
}

func (m *myStore) Remove(ctx context.Context, id string) error {
	m.c.Delete(id)
	return nil
}

func (m *myStore) Get(ctx context.Context, id string) (Session, error) {
	res, ok := m.c.Get(id)
	if !ok {
		return nil, fmt.Errorf("session 不存在，id：%s", id)
	}
	return res.(Session), nil
}

// Session 基于内存的实现
// 需要加锁，因为有可能同一个用户会并发发起请求
type memorySession struct {
	data  map[string]string
	id    string
	mutex sync.RWMutex
}

func (s *memorySession) Get(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("key 不存在")
	}
	return v, nil
}

func (s *memorySession) Set(ctx context.Context, key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
	return nil
}

func (s *memorySession) ID() string {
	return s.id
}
