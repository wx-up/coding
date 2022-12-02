package pool

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// ControlPool 装饰器增强 sync.Pool
// 在对 sync.Pool 进行封装的时候，需要注意它是被 GC 托管的
type ControlPool struct {
	pool sync.Pool

	maxCnt int32
	cnt    int32
}

func (p *ControlPool) Get() any {
	return p.pool.Get()
}

// Put 放入缓存
// 当前实现存在问题，因为 sync.Pool 托管给 GC，实际 Pool 中的对象数并不是 cnt 统计的数量
// 由于 GC 的存在会绕过自己设计的控制
func (p *ControlPool) Put(val any) {
	// 单个缓存对象超过 1024 个字节，直接 return 不放回 pool
	if unsafe.Sizeof(val) > 1024 {
		return
	}

	// 先占个坑
	cnt := atomic.AddInt32(&p.cnt, 1)

	// 判断是否超过最大个数的限制
	if cnt > p.maxCnt {
		atomic.AddInt32(&p.cnt, -1)
		return
	}

	// 放入 pool
	p.pool.Put(val)
}
