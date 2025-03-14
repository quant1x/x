package ringbuffer

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	ErrMaxConsumers   = errors.New("max consumers reached")
	ErrInvalidSize    = errors.New("size must be power of two")
	ErrConsumerClosed = errors.New("consumer closed")
)

// 队列核心结构
type RingBuffer[T any] struct {
	buffer      []T           // 数据存储
	size        uint32        // 队列容量（2^n）
	mask        uint32        // 位掩码（size-1）
	producerPos uint32        // 生产者位置（原子操作）
	consumers   []consumerPtr // 消费者指针数组
}

// 消费者指针（内存对齐优化）
type consumerPtr struct {
	posVersion uint64   // 高32位:版本号, 低32位:位置
	active     uint32   // 活跃状态（0: 未激活, 1: 激活）
	_          [52]byte // 填充至64字节（避免伪共享）
}

// 消费者句柄
type Consumer[T any] struct {
	rb    *RingBuffer[T]
	index uint32 // 在consumers数组中的索引
}

// 创建队列
func New[T any](size uint32, maxConsumers int) (*RingBuffer[T], error) {
	if size&(size-1) != 0 {
		return nil, ErrInvalidSize
	}

	rb := &RingBuffer[T]{
		buffer:    make([]T, size),
		size:      size,
		mask:      size - 1,
		consumers: make([]consumerPtr, maxConsumers),
	}

	// 强制内存对齐（确保posVersion为64位对齐）
	if uintptr(unsafe.Pointer(&rb.consumers[0]))%8 != 0 {
		return nil, errors.New("consumerPtr not aligned")
	}

	return rb, nil
}

// 创建新消费者
func (rb *RingBuffer[T]) NewConsumer() (*Consumer[T], error) {
	initialPos := atomic.LoadUint32(&rb.producerPos)

	for i := range rb.consumers {
		ptr := &rb.consumers[i]
		// 使用CAS激活消费者
		if atomic.CompareAndSwapUint32(&ptr.active, 0, 1) {
			// 初始化位置和版本号 (version=0, pos=initialPos)
			atomic.StoreUint64(&ptr.posVersion, uint64(initialPos))
			return &Consumer[T]{rb: rb, index: uint32(i)}, nil
		}
	}
	return nil, ErrMaxConsumers
}

// 关闭消费者
func (c *Consumer[T]) Close() {
	ptr := &c.rb.consumers[c.index]
	atomic.StoreUint32(&ptr.active, 0)
}

// 生产者写入
func (rb *RingBuffer[T]) Write(value T) {
	for {
		currentProd := atomic.LoadUint32(&rb.producerPos)
		minPos := rb.findMinConsumerPos()

		// 队列已满检查
		if currentProd-minPos >= rb.size {
			runtime.Gosched()
			continue
		}

		// 预占写入位置
		newProd := currentProd + 1
		if atomic.CompareAndSwapUint32(&rb.producerPos, currentProd, newProd) {
			// 写入数据
			rb.buffer[currentProd&rb.mask] = value
			return
		}
	}
}

// 消费者读取（修复ABA问题）
func (c *Consumer[T]) Read() (T, error) {
	ptr := &c.rb.consumers[c.index]

	for {
		// 1. 原子加载位置和版本号
		current := atomic.LoadUint64(&ptr.posVersion)
		currentVersion := uint32(current >> 32) // 高32位为版本号
		currentPos := uint32(current)           // 低32位为位置

		// 2. 检查消费者是否关闭
		if atomic.LoadUint32(&ptr.active) == 0 {
			var zero T
			return zero, ErrConsumerClosed
		}

		// 3. 检查数据是否就绪
		currentProd := atomic.LoadUint32(&c.rb.producerPos)
		if currentPos >= currentProd {
			runtime.Gosched()
			continue
		}

		// 4. 读取数据
		index := currentPos & c.rb.mask
		value := c.rb.buffer[index]

		// 5. 准备新位置和版本号
		newPos := currentPos + 1
		newVersion := currentVersion + 1 // 版本号递增
		new := (uint64(newVersion) << 32) | uint64(newPos)

		// 6. CAS更新（同时校验位置和版本号）
		if atomic.CompareAndSwapUint64(&ptr.posVersion, current, new) {
			return value, nil
		}

		// 7. 如果CAS失败，重试循环
	}
}

// 查找最慢消费者位置
func (rb *RingBuffer[T]) findMinConsumerPos() uint32 {
	minPos := atomic.LoadUint32(&rb.producerPos)
	for i := range rb.consumers {
		ptr := &rb.consumers[i]
		if atomic.LoadUint32(&ptr.active) == 1 {
			// 从posVersion中提取位置
			current := atomic.LoadUint64(&ptr.posVersion)
			currentPos := uint32(current)
			if currentPos < minPos {
				minPos = currentPos
			}
		}
	}
	return minPos
}
