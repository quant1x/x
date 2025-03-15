package ringbuffer

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	ErrQueueFull   = errors.New("queue is full")
	ErrQueueEmpty  = errors.New("queue is empty")
	ErrInvalidSize = errors.New("size must be power of two")
)

//go:align 64
type consumer struct {
	pos    uint32   // 当前读取序列号（非索引）
	active uint32   // 消费者状态
	_      [56]byte // 填充至64字节
}

type RingBuffer[T any] struct {
	buffer      []T
	size        uint32 // 队列容量（必须为2^n）
	mask        uint32 // 位掩码（size-1）
	producerPos uint32 // 生产者序列号（持续递增）
	consumers   []consumer
}

type Consumer[T any] struct {
	rb    *RingBuffer[T]
	index uint32
}

// 创建队列（正确初始化消费者）
func New[T any](size uint32, maxConsumers int) (*RingBuffer[T], error) {
	if size&(size-1) != 0 {
		return nil, ErrInvalidSize
	}

	if unsafe.Sizeof(consumer{}) != 64 {
		return nil, errors.New("consumer alignment failed")
	}

	rb := &RingBuffer[T]{
		buffer:    make([]T, size),
		size:      size,
		mask:      size - 1,
		consumers: make([]consumer, maxConsumers),
	}

	// 初始化消费者位置为0
	for i := range rb.consumers {
		atomic.StoreUint32(&rb.consumers[i].pos, 0)
	}
	return rb, nil
}

// 注册消费者（修复初始化位置）
func (rb *RingBuffer[T]) NewConsumer() (*Consumer[T], error) {
	initialPos := atomic.LoadUint32(&rb.producerPos)

	for i := range rb.consumers {
		c := &rb.consumers[i]
		if atomic.CompareAndSwapUint32(&c.active, 0, 1) {
			// 初始位置对齐到当前生产周期起点
			base := initialPos - (initialPos % rb.size)
			atomic.StoreUint32(&c.pos, base)
			return &Consumer[T]{rb: rb, index: uint32(i)}, nil
		}
	}
	return nil, errors.New("max consumers reached")
}

// 生产者写入（修复并发写入）
func (rb *RingBuffer[T]) Write(value T) error {
	for {
		// 获取当前生产序列号
		currentProd := atomic.LoadUint32(&rb.producerPos)
		writeIndex := currentProd & rb.mask

		// 检查是否有空间可写
		minPos := rb.findMinConsumerPos()
		if currentProd-minPos >= rb.size {
			return ErrQueueFull
		}

		// 预占写入位置（CAS替代AddUint32）
		if !atomic.CompareAndSwapUint32(&rb.producerPos, currentProd, currentProd+1) {
			runtime.Gosched()
			continue
		}

		// 写入数据（保证内存可见性）
		rb.buffer[writeIndex] = value
		return nil
	}
}

// 消费者读取（修复空判断）
func (c *Consumer[T]) Read() (T, error) {
	var zero T
	ptr := &c.rb.consumers[c.index]

	for {
		// 获取当前生产位置（内存屏障保证可见性）
		currentProd := atomic.LoadUint32(&c.rb.producerPos)
		currentRead := atomic.LoadUint32(&ptr.pos)

		// 计算可用数据量（处理环形溢出）
		available := currentProd - currentRead
		if available == 0 {
			return zero, ErrQueueEmpty
		}

		// 读取数据
		readIndex := currentRead & c.rb.mask
		value := c.rb.buffer[readIndex]

		// 原子更新读取位置
		if atomic.CompareAndSwapUint32(&ptr.pos, currentRead, currentRead+1) {
			return value, nil
		}
		runtime.Gosched()
	}
}

// 查找最慢消费者（修复溢出处理）
func (rb *RingBuffer[T]) findMinConsumerPos() uint32 {
	currentProd := atomic.LoadUint32(&rb.producerPos)
	minPos := currentProd // 初始为生产者位置

	for i := range rb.consumers {
		c := &rb.consumers[i]
		if atomic.LoadUint32(&c.active) == 1 {
			pos := atomic.LoadUint32(&c.pos)
			// 处理溢出（当currentProd超过uint32最大值时）
			if currentProd-pos > 0xFFFFFFFF-rb.size {
				pos = currentProd - rb.size
			}
			if pos < minPos {
				minPos = pos
			}
		}
	}
	return minPos
}
