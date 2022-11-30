package decorator

import (
	"errors"
	"sync"
)

type Cache interface {
	Get(string) (string, error)
}

type memoryCache struct {
	values map[string]string
}

func (m *memoryCache) Get(key string) (string, error) {
	v, ok := m.values[key]
	if !ok {
		return "", errors.New("key not exist")
	}
	return v, nil
}

// SafeCache 对已有功能进行增强（ 接口不会改变 ）
type SafeCache struct {
	Cache
	lock sync.RWMutex
}

func (sc *SafeCache) Get(key string) (string, error) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	return sc.Cache.Get(key)
}
