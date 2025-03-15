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
	pos    uint32
	active uint32
	_      [56]byte // 填充至64字节
}

type RingBuffer[T any] struct {
	buffer      []atomic.Value // 使用atomic.Value存储数据
	size        uint32
	mask        uint32
	producerPos uint32
	consumers   []consumer
}

type Consumer[T any] struct {
	rb    *RingBuffer[T]
	index uint32
}

// 创建队列
func New[T any](size uint32, maxConsumers int) (*RingBuffer[T], error) {
	if size&(size-1) != 0 {
		return nil, ErrInvalidSize
	}

	rb := &RingBuffer[T]{
		buffer:    make([]atomic.Value, size),
		size:      size,
		mask:      size - 1,
		consumers: make([]consumer, maxConsumers),
	}

	if unsafe.Sizeof(consumer{}) != 64 {
		return nil, errors.New("consumer alignment failed")
	}
	return rb, nil
}

// 注册消费者
func (rb *RingBuffer[T]) NewConsumer() (*Consumer[T], error) {
	initialPos := atomic.LoadUint32(&rb.producerPos)

	for i := range rb.consumers {
		c := &rb.consumers[i]
		if atomic.CompareAndSwapUint32(&c.active, 0, 1) {
			base := initialPos - (initialPos % rb.size)
			atomic.StoreUint32(&c.pos, base)
			return &Consumer[T]{rb: rb, index: uint32(i)}, nil
		}
	}
	return nil, errors.New("max consumers reached")
}

// 生产者写入（原子存储）
func (rb *RingBuffer[T]) Write(value T) error {
	for {
		currentProd := atomic.LoadUint32(&rb.producerPos)
		minPos := rb.findMinConsumerPos()

		if currentProd-minPos >= rb.size {
			return ErrQueueFull
		}

		newProd := currentProd + 1
		if !atomic.CompareAndSwapUint32(&rb.producerPos, currentProd, newProd) {
			runtime.Gosched()
			continue
		}

		// 原子存储数据
		index := currentProd & rb.mask
		rb.buffer[index].Store(value)
		return nil
	}
}

// 消费者读取（原子加载）
func (c *Consumer[T]) Read() (T, error) {
	var zero T
	ptr := &c.rb.consumers[c.index]

	for {
		currentProd := atomic.LoadUint32(&c.rb.producerPos)
		currentRead := atomic.LoadUint32(&ptr.pos)

		if currentRead >= currentProd {
			return zero, ErrQueueEmpty
		}

		// 原子加载数据
		readIndex := currentRead & c.rb.mask
		val := c.rb.buffer[readIndex].Load()
		if val == nil {
			runtime.Gosched()
			continue
		}

		if atomic.CompareAndSwapUint32(&ptr.pos, currentRead, currentRead+1) {
			return val.(T), nil
		}
		runtime.Gosched()
	}
}

// 查找最慢消费者（逻辑不变）
func (rb *RingBuffer[T]) findMinConsumerPos() uint32 {
	currentProd := atomic.LoadUint32(&rb.producerPos)
	minPos := currentProd

	for i := range rb.consumers {
		c := &rb.consumers[i]
		if atomic.LoadUint32(&c.active) == 1 {
			pos := atomic.LoadUint32(&c.pos)
			if currentProd-pos >= rb.size {
				pos = currentProd - rb.size
			}
			if pos < minPos {
				minPos = pos
			}
		}
	}
	return minPos
}
