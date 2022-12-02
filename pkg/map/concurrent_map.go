package _map

import "sync"

// SafeMap 并发安全的 map
// comparable 表示可比较的，约束 K，比如 slice 就不能被当作 K 因为它是不可比较的
type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

// LoadOrStore k 存在时返回 v，k 不存在是设置并返回v
func (sm *SafeMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	sm.lock.RLock()
	val, ok := sm.values[k]
	// 不能使用 defer 释放，否则会导致死锁
	sm.lock.RUnlock()
	if ok {
		return val, true
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()

	// double check
	val, ok = sm.values[k]
	if ok {
		return val, true
	}

	sm.values[k] = v
	return v, false
}

// LoadOrStoreHeavy 函数式编程，延迟初始化
func (sm *SafeMap[K, V]) LoadOrStoreHeavy(k K, f func() V) (V, bool) {
	sm.lock.RLock()
	val, ok := sm.values[k]
	sm.lock.RUnlock()
	if ok {
		return val, true
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()

	val, ok = sm.values[k]
	if ok {
		return val, true
	}

	// 延迟初始化
	newVal := f()

	sm.values[k] = newVal
	return newVal, false
}
