package ringbuffer

import (
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestMPMCRingBuffer(t *testing.T) {
	t.Parallel()

	const size = 1024
	rb, _ := New[int](size)

	numProducers := 4
	numConsumers := 4
	dataPerProducer := 10000
	totalData := numProducers * dataPerProducer

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	collected := make(chan int, totalData)

	// 生产者
	producerWg.Add(numProducers)
	for i := 0; i < numProducers; i++ {
		go func(id int) {
			defer producerWg.Done()
			for j := 0; j < dataPerProducer; j++ {
				value := id*dataPerProducer + j
				for {
					if err := rb.Write(value); err == nil {
						break
					}
					runtime.Gosched()
				}
			}
		}(i)
	}

	// 消费者
	consumerWg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go func() {
			defer consumerWg.Done()
			for {
				v, err := rb.Read()
				if err != nil {
					break // 包括 ErrClosed
				}
				collected <- v
			}
		}()
	}

	// 关闭：生产者完成 -> Close
	go func() {
		producerWg.Wait()
		rb.Close()
	}()

	// 收集完成
	go func() {
		consumerWg.Wait()
		close(collected)
	}()

	// 收集结果
	result := make([]int, 0, totalData)
	for v := range collected {
		result = append(result, v)
	}

	// 验证数量
	if len(result) != totalData {
		t.Fatalf("expected %d items, got %d", totalData, len(result))
	}

	// 验证内容（排序后）
	sort.Ints(result)
	for i := 0; i < totalData; i++ {
		if result[i] != i {
			t.Fatalf("result[%d] = %d, want %d", i, result[i], i)
		}
	}
}

// 压力测试
func TestStressMPMCRingBuffer(t *testing.T) {
	t.Parallel()
	start := time.Now()

	const size = 65536 // 64K slots
	rb, err := New[int](size)
	if err != nil {
		t.Fatal(err)
	}

	numProducers := 8
	numConsumers := 8
	dataPerProducer := 30000
	totalData := numProducers * dataPerProducer

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// 使用无缓冲 channel 收集，避免内存膨胀
	collected := make(chan int, totalData)

	// 生产者：高并发写入
	producerWg.Add(numProducers)
	for i := 0; i < numProducers; i++ {
		go func(id int) {
			defer producerWg.Done()
			for j := 0; j < dataPerProducer; j++ {
				value := id*dataPerProducer + j
				for {
					if err := rb.Write(value); err == nil {
						break
					}
					// 只有 ErrQueueFull 才重试，其他错误 panic
					if err != ErrQueueFull {
						t.Errorf("unexpected write error: %v", err)
						return
					}
					runtime.Gosched()
				}
			}
		}(i)
	}

	// 消费者：读取并发送到 channel
	consumerWg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go func() {
			defer consumerWg.Done()
			for {
				v, err := rb.Read()
				if err != nil {
					// ErrClosed 是正常退出
					return
				}
				_ = v
				select {
				case collected <- v:
				default:
					// 防止 collected 阻塞
				}
			}
		}()
	}

	// 关闭逻辑：生产者完成 → Close
	go func() {
		producerWg.Wait()
		rb.Close()
	}()

	// 收集结果
	go func() {
		consumerWg.Wait()
		close(collected)
	}()

	// 验证数据
	seen := make(map[int]bool, totalData)
	var count int
	for v := range collected {
		if v < 0 || v >= totalData {
			t.Errorf("out of range value: %d", v)
		}
		if seen[v] {
			t.Errorf("duplicate value: %d", v)
		}
		seen[v] = true
		count++
	}

	if count != totalData {
		t.Fatalf("expected %d items, got %d", totalData, count)
	}

	if len(seen) != totalData {
		t.Fatalf("missing or duplicate data: got %d unique", len(seen))
	}

	//producerWg.Wait()
	//rb.Close()
	//consumerWg.Wait()
	t.Logf("Stress test passed: %d producers, %d consumers, %d items", numProducers, numConsumers, totalData)
	t.Logf("cross time:%v", time.Since(start))
}

// 基准测试
func BenchmarkRingBuffer_WriteRead(b *testing.B) {
	rb, _ := New[int](65536)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// 随机选择写或读（基于 goroutine 的地址或简单哈希）
		// 这里我们用一个简单方式：根据 goroutine ID 的奇偶性决定角色
		// 实际上，我们让所有 goroutine 都做写和读的混合操作

		for pb.Next() {
			// 50% 概率写，50% 概率读
			if b.N%2 == 0 || b.N%3 == 0 {
				// 写入
				for {
					if rb.Write(42) == nil {
						break
					}
					runtime.Gosched()
				}
			} else {
				// 读取
				_, err := rb.Read()
				if err != nil {
					runtime.Gosched()
				}
			}
		}
	})
}

// 对比基准
func BenchmarkRingBuffer_WriteOnly(b *testing.B) {
	rb, _ := New[int](65536)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				if rb.Write(42) == nil {
					break
				}
				runtime.Gosched()
			}
		}
	})
}

func BenchmarkRingBuffer_ReadOnly(b *testing.B) {
	rb, _ := New[int](65536)
	// 预写入一些数据
	for i := 0; i < 1000; i++ {
		rb.Write(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := rb.Read()
			if err != nil {
				// 如果空了，重写一些
				rb.Write(42)
			}
		}
	})
}

func BenchmarkRingBuffer_Mixed(b *testing.B) {
	rb, _ := New[int](65536)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 混合读写
		for pb.Next() {
			if b.N%2 == 0 {
				// 写
				for {
					if rb.Write(42) == nil {
						break
					}
					runtime.Gosched()
				}
			} else {
				// 读
				_, err := rb.Read()
				if err != nil {
					runtime.Gosched()
				}
			}
		}
	})
}

// 验证关闭行为
func TestRingBuffer_CloseBehavior(t *testing.T) {
	rb, _ := New[string](16)

	// 写入一些数据
	err := rb.Write("hello")
	if err != nil {
		t.Fatal(err)
	}

	rb.Close()

	// 读取已写入的数据
	v, err := rb.Read()
	if err != nil {
		t.Fatalf("should read existing data, but got error: %v", err)
	}
	if v != "hello" {
		t.Fatalf("expected 'hello', got %q", v)
	}

	// 再读一次，应该返回 ErrClosed
	_, err = rb.Read()
	if err == nil {
		t.Fatal("expected error after close and empty, got nil")
	}
	if err != ErrClosed {
		t.Fatalf("expected ErrClosed, got %v", err)
	}

	// 写入应失败
	err = rb.Write("world")
	if err == nil {
		t.Fatal("expected write to fail after close")
	}
}
