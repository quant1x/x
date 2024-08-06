package main

import (
	"fmt"
	"sync"
)

// SafeMap 定义一个并发安全的 map
type SafeMap struct {
	mu sync.Mutex
	m  map[string]int
}

// NewSafeMap 创建一个新的并发安全的 map
func NewSafeMap() *SafeMap {

	return &SafeMap{

		m: make(map[string]int),
	}
}

// 设置键值对，加锁保护
func (s *SafeMap) Set(key string, value int) {

	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// Get 根据键获取值，加锁保护
func (s *SafeMap) Get(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.m[key]
	return val, ok
}

// go run -gcflags "-m" github.com/quant1x/x/examples/escape-analysis
func main() {
	sm := NewSafeMap()
	// 设置值
	sm.Set("hello", 42)
	// 获取值
	var val int
	var ok bool
	if val, ok = sm.Get("hello"); ok {
		fmt.Printf("Value: %d\n", val)
	}
	sm = nil
}
