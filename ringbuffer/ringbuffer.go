package ringbuffer

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	ErrQueueFull   = errors.New("queue is full")
	ErrInvalidSize = errors.New("size must be power of two")
	ErrClosed      = errors.New("ring buffer closed")
)

// Slot represents a single slot in the ring buffer
type Slot[T any] struct {
	data unsafe.Pointer // 数据存储
	flag uint32         // 状态标志 (0: empty, 1: writing, 2: readable)
}

// RingBuffer represents the MPMC ring buffer
type RingBuffer[T any] struct {
	slots       []Slot[T] // 使用槽位数组存储数据
	size        uint32
	mask        uint32
	producerPos uint32 // 全局生产者位置
	consumerPos uint32 // 全局消费者位置
	closed      uint32 // 关闭标记
}

// New creates a new MPMC ring buffer
func New[T any](size uint32) (*RingBuffer[T], error) {
	if size == 0 || (size&(size-1)) != 0 {
		return nil, ErrInvalidSize
	}

	rb := &RingBuffer[T]{
		slots: make([]Slot[T], size),
		size:  size,
		mask:  size - 1,
	}

	for i := range rb.slots {
		atomic.StoreUint32(&rb.slots[i].flag, 0) // 初始化为empty状态
	}

	return rb, nil
}

// spinWait 自旋等待，使用指数退避
func spinWait(retries *int32) {
	r := atomic.AddInt32(retries, 1)
	switch {
	case r < 4:
		for i := 0; i < 1<<(r*2); i++ {
		}
	case r < 8:
		runtime.Gosched()
	default:
		runtime.Gosched()
	}
}

// Write writes data into the ring buffer by a producer.
//
// It explicitly boxes the value on the heap to ensure the pointer remains valid,
// avoiding reliance on compiler escape analysis for safety.
func (rb *RingBuffer[T]) Write(value T) error {
	if atomic.LoadUint32(&rb.closed) == 1 {
		return errors.New("queue closed")
	}

	var currentProd, minCons uint32
	var retries int32 = 0 // ✅ 引入重试计数
	for {
		currentProd = atomic.LoadUint32(&rb.producerPos)
		minCons = atomic.LoadUint32(&rb.consumerPos)

		if currentProd-minCons >= rb.size {
			spinWait(&retries)
			continue
		}

		index := currentProd & rb.mask
		slot := &rb.slots[index]

		// 尝试获取写权限
		if atomic.LoadUint32(&slot.flag) != 0 {
			spinWait(&retries)
			continue
		}

		// CAS更新槽位状态为writing
		if !atomic.CompareAndSwapUint32(&slot.flag, 0, 1) {
			spinWait(&retries)
			continue
		}

		// 写入数据：显式堆分配，不再依赖逃逸分析
		boxed := new(T)
		*boxed = value
		atomic.StorePointer(&slot.data, unsafe.Pointer(boxed))
		atomic.StoreUint32(&slot.flag, 2)

		// 更新全局生产者位置
		if atomic.CompareAndSwapUint32(&rb.producerPos, currentProd, currentProd+1) {
			return nil
		}

		// 如果更新失败，回滚槽位状态
		atomic.StoreUint32(&slot.flag, 0)
		spinWait(&retries)
	}
}

// Read reads data from the ring buffer by a consumer
func (rb *RingBuffer[T]) Read() (T, error) {
	var zero T
	var retries int32 = 0 // ✅ 引入重试计数
	for {
		currentCons := atomic.LoadUint32(&rb.consumerPos)
		currentProd := atomic.LoadUint32(&rb.producerPos)

		// 检查队列是否关闭且没有更多数据
		if atomic.LoadUint32(&rb.closed) == 1 && currentCons >= currentProd {
			return zero, ErrClosed
		}

		if currentCons >= currentProd {
			// 队列为空时让出 CPU 时间片
			spinWait(&retries)
			continue
		}

		index := currentCons & rb.mask
		slot := &rb.slots[index]

		// 检查槽位是否可读
		if atomic.LoadUint32(&slot.flag) != 2 {
			spinWait(&retries)
			continue
		}

		// CAS更新槽位状态为empty
		if !atomic.CompareAndSwapUint32(&slot.flag, 2, 0) {
			spinWait(&retries)
			continue
		}

		// 读取数据并更新全局消费者位置
		valPtr := atomic.LoadPointer(&slot.data)
		if valPtr == nil {
			atomic.StoreUint32(&slot.flag, 2) // 回滚槽位状态
			spinWait(&retries)
			continue
		}

		if atomic.CompareAndSwapUint32(&rb.consumerPos, currentCons, currentCons+1) {
			return *(*T)(valPtr), nil
		}

		// 如果更新失败，回滚槽位状态
		atomic.StoreUint32(&slot.flag, 2)
		spinWait(&retries)
	}
}

// Close closes the ring buffer
func (rb *RingBuffer[T]) Close() {
	atomic.StoreUint32(&rb.closed, 1)
}
