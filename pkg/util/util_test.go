package util

import "testing"

// 1. 基本的for循环测试
func BenchmarkBasicLoop(b *testing.B) {
	arr := make([]int, 1000)
	for i := 0; i < b.N; i++ {
		for j := range arr {
			_ = j
		}
	}
}

// 2. 测试数组遍历
func BenchmarkArrayLoop(b *testing.B) {
	arr := make([]int, 1000)
	b.ResetTimer() // 重置计时器，去掉初始化影响
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(arr); j++ {
			_ = arr[j]
		}
	}
}
