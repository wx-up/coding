package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := make([]int, 10)
	fmt.Println(len(s))
	fmt.Println(cap(s))
}

func Add(a, b int) int {
	return a / b
}

// FuzzAdd 模糊测试
// 如果要基于种子语料库生成随机测试数据用于模糊测试，需要给go test命令增加 -fuzz参数
//
//	go test -fuzz=Fuzz
//
// 测试失败后，失败的用例会输出到文件 testdata/fuzz/FuzzXXX/ 目录下
// 模糊测试和单元测试是互补的，谁也替代不了谁
func FuzzAdd(f *testing.F) {
	// 第一部分是添加初始的测试输入，这些测试输入可以被看做是种子数据或样本数据，也被称为 Seed Corpus
	f.Add(1, 100)

	// 第二部分是执行目标测试函数
	f.Fuzz(func(t *testing.T, a int, b int) {
		Add(a, b)
	})
}
