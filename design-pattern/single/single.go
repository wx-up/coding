package single

import (
	"fmt"
	"sync"
)

// singleton 单例
// 不可导出，外部不能直接创建实例
// 暴露一个构建方法，返回实例，比如：GetSingleInstance
type singleton struct{}

func (s *singleton) Do() {
	fmt.Println("我是单例")
}

var (
	instance    *singleton
	instanceOne sync.Once
)

func GetSingleInstance() *singleton {
	// 只会执行一次
	instanceOne.Do(func() {
		instance = &singleton{}
	})
	return instance
}
