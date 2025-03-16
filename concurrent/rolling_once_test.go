package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewDailyOnce(t *testing.T) {
	cache := NewDailyOnce()

	// 第一次调用触发初始化
	cache.Do(func() {
		fmt.Println("执行初始化 @", time.Now().Format("15:04:05"))
	})

	// 当天后续调用直接返回
	cache.Do(func() {
		fmt.Println("这行不应该出现")
	})

	// 模拟第二天调用（调整系统时间测试）
	fmt.Println("\n模拟第二天调用:")
	cache.Do(func() {
		fmt.Println("执行新一天初始化 @", time.Now().Format("2006-01-02 15:04:05"))
	})
}

func BenchmarkStandardOnce(b *testing.B) {
	o := &sync.Once{}
	for i := 0; i < b.N; i++ {
		o.Do(func() {})
	}
}

func BenchmarkDailyOnce(b *testing.B) {
	d := NewDailyOnce()
	for i := 0; i < b.N; i++ {
		d.Do(func() {})
	}
}

func BenchmarkPeriodOnce(b *testing.B) {
	d := CreatePeriodOnceWithSecond(5)
	for i := 0; i < b.N; i++ {
		d.Do(func() {})
	}
}
