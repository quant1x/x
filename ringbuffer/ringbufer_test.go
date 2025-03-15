package ringbuffer

import (
	"reflect"
	"runtime"
	"sort"
	"sync"
	"testing"
)

func TestMPMCRingBuffer(t *testing.T) {
	const size = 1024
	rb, _ := New[int](size)

	numProducers := 2
	numConsumers := 2
	dataPerProducer := 5000
	totalData := numProducers * dataPerProducer

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)

	collected := make(chan int, totalData) // 用于收集消费者读取的数据

	// 生产者协程
	for i := 0; i < numProducers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < dataPerProducer; j++ {
				value := id*dataPerProducer + j
				for {
					err := rb.Write(value)
					if err == nil {
						break
					}
					if err == ErrQueueFull {
						runtime.Gosched()
						continue
					}
					t.Fatalf("unexpected error during write: %v", err)
				}
			}
		}(i)
	}

	// 消费者协程
	for i := 0; i < numConsumers; i++ {
		go func() {
			defer wg.Done()
			for {
				v, err := rb.Read()
				if err != nil {
					break
				}
				collected <- v
			}
		}()
	}

	// 等待所有生产者完成写入
	go func() {
		wg.Wait() // 等待所有生产者和消费者完成
		rb.Close()
	}()

	// 收集所有消费者读取的数据
	result := make([]int, 0, totalData)
	for v := range collected {
		result = append(result, v)
	}

	// 验证数据完整性
	sort.Ints(result)
	expected := make([]int, totalData)
	for i := 0; i < totalData; i++ {
		expected[i] = i
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("data mismatch:\ncollected: %v\nexpected: %v", result, expected)
	}
}
