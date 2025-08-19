package ringbuffer

import (
	"runtime"
	"sort"
	"sync"
	"testing"
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
