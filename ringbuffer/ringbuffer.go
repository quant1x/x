package ringbuffer

import (
	"errors"
	"math"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	ErrQueueFull   = errors.New("queue is full")
	ErrInvalidSize = errors.New("size must be power of two")
)

//go:align 64
type consumer struct {
	pos    uint32
	active uint32
	_      [56]byte // 填充至64字节，避免伪共享
}

type RingBuffer[T any] struct {
	buffer      []unsafe.Pointer // 使用unsafe.Pointer存储数据
	size        uint32
	mask        uint32
	producerPos uint32
	consumers   []consumer
	closed      uint32
}

type Consumer[T any] struct {
	rb    *RingBuffer[T]
	index uint32
}

func New[T any](size uint32, maxConsumers int) (*RingBuffer[T], error) {
	if size == 0 || (size&(size-1)) != 0 {
		return nil, ErrInvalidSize
	}

	rb := &RingBuffer[T]{
		buffer:    make([]unsafe.Pointer, size),
		size:      size,
		mask:      size - 1,
		consumers: make([]consumer, maxConsumers),
	}

	if unsafe.Sizeof(consumer{}) != 64 {
		return nil, errors.New("consumer alignment failed")
	}

	for i := range rb.buffer {
		atomic.StorePointer(&rb.buffer[i], nil)
	}

	return rb, nil
}

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

func (rb *RingBuffer[T]) Write(value T) error {
	if atomic.LoadUint32(&rb.closed) == 1 {
		return errors.New("queue closed")
	}

	var currentProd, minPos uint32
	for {
		currentProd = atomic.LoadUint32(&rb.producerPos)
		minPos = rb.findMinConsumerPos()

		if currentProd-minPos >= rb.size {
			return ErrQueueFull
		}

		valPtr := unsafe.Pointer(&value)
		index := currentProd & rb.mask

		// 写入内存屏障确保数据可见性
		atomic.StorePointer(&rb.buffer[index], valPtr)
		if atomic.CompareAndSwapUint32(&rb.producerPos, currentProd, currentProd+1) {
			return nil
		}
		runtime.Gosched()
	}
}

func (c *Consumer[T]) Read() (T, error) {
	var zero T
	ptr := &c.rb.consumers[c.index]

	for {
		currentProd := atomic.LoadUint32(&c.rb.producerPos)
		currentRead := atomic.LoadUint32(&ptr.pos)

		// 检查队列是否关闭
		if atomic.LoadUint32(&c.rb.closed) == 1 {
			return zero, errors.New("queue closed")
		}

		if currentRead >= currentProd {
			// 队列为空时让出 CPU 时间片
			runtime.Gosched()
			continue
		}

		readIndex := currentRead & c.rb.mask
		valPtr := atomic.LoadPointer(&c.rb.buffer[readIndex])

		if valPtr == nil {
			runtime.Gosched()
			continue
		}

		// 双重检查确保数据有效性
		const maxRetries = 1000
		retries := 0
		for {
			if atomic.CompareAndSwapUint32(&ptr.pos, currentRead, currentRead+1) {
				break
			}
			retries++
			if retries > maxRetries {
				return zero, errors.New("consumer position update failed after max retries")
			}
			runtime.Gosched()
		}

		// 确保读取数据后清空槽位
		atomic.StorePointer(&c.rb.buffer[readIndex], nil)
		return *(*T)(valPtr), nil
	}
}

func (rb *RingBuffer[T]) findMinConsumerPos() uint32 {
	minPos := uint32(math.MaxUint32)
	for i := range rb.consumers {
		c := &rb.consumers[i]
		if atomic.LoadUint32(&c.active) == 1 {
			pos := atomic.LoadUint32(&c.pos)
			if pos < minPos {
				minPos = pos
			}
		}
	}
	if minPos == math.MaxUint32 {
		return atomic.LoadUint32(&rb.producerPos)
	}
	return minPos
}

func (rb *RingBuffer[T]) Close() {
	atomic.StoreUint32(&rb.closed, 1)
}
