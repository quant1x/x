package main

import (
	"fmt"
	"github.com/quant1x/x/ringbuffer"
	"time"
)

func main() {
	count := 1000
	consumerNum := 2
	rb := ringbuffer.NewRingBuffer(4) // 创建容量为4的队列

	// 监控队列状态
	go func() {
		for {
			fmt.Println("IsEmpty:", rb.IsEmpty(), "IsFull:", rb.IsFull())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 消费者
	for i := 0; i < consumerNum; i++ {
		go func(no int) {
			for {
				if val, ok := rb.Dequeue(); ok {
					fmt.Printf("No%d: Dequeued: %v\n", no, val)
				}
				//time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}
	time.Sleep(time.Second)
	// 生产者
	go func() {
		for i := 0; i < count; i++ {
			if rb.Enqueue(i) {
				fmt.Printf("Enqueued: %d\n", i)
			} else {
				fmt.Printf("Enqueue failed: %d\n", i)
			}
			//time.Sleep(200 * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)
}
