package ringbuffer

import (
	"sync/atomic"
)

type RingBuffer struct {
	data []atomic.Value // 使用atomic.Value切片存储数据
	size int32          // 队列容量，必须为2的幂
	mask int32          // 用于快速取模运算
	tail int32          // 生产者指针（下一个写入位置）
	head int32          // 消费者指针（下一个读取位置）
}

func NewRingBuffer(size int32) *RingBuffer {
	if size&(size-1) != 0 {
		panic("size must be a power of two")
	}
	return &RingBuffer{
		data: make([]atomic.Value, size),
		size: size,
		mask: size - 1,
	}
}

// IsEmpty 判断队列是否为空（线程安全）
func (rb *RingBuffer) IsEmpty() bool {
	// 注意加载顺序：先读head再读tail，与消费者逻辑一致
	currentHead := atomic.LoadInt32(&rb.head)
	currentTail := atomic.LoadInt32(&rb.tail)
	return currentHead == currentTail
}

// IsFull 判断队列是否已满（线程安全）
func (rb *RingBuffer) IsFull() bool {
	// 注意加载顺序：先读tail再读head，与生产者逻辑一致
	currentTail := atomic.LoadInt32(&rb.tail)
	currentHead := atomic.LoadInt32(&rb.head)
	return currentTail-currentHead >= rb.size
}

// Enqueue 添加元素到队列，成功返回true，队列满时返回false
func (rb *RingBuffer) Enqueue(value interface{}) bool {
	for {
		currentTail := atomic.LoadInt32(&rb.tail)
		currentHead := atomic.LoadInt32(&rb.head)
		if currentTail-currentHead >= rb.size { // 队列已满
			return false
		}

		// 尝试抢占当前tail位置
		if atomic.CompareAndSwapInt32(&rb.tail, currentTail, currentTail+1) {
			index := currentTail & rb.mask
			rb.data[index].Store(value) // 原子存储数据
			return true
		}
	}
}

// Dequeue 从队列取出元素，成功返回元素和true，队列空时返回nil和false
func (rb *RingBuffer) Dequeue() (interface{}, bool) {
	for {
		currentHead := atomic.LoadInt32(&rb.head)
		currentTail := atomic.LoadInt32(&rb.tail)
		if currentHead == currentTail { // 队列为空
			return nil, false
		}

		index := currentHead & rb.mask
		value := rb.data[index].Load() // 原子读取数据

		if value == nil { // 数据未就绪时重试
			continue
		}

		// 尝试抢占当前head位置
		if atomic.CompareAndSwapInt32(&rb.head, currentHead, currentHead+1) {
			return value, true
		}
	}
}
