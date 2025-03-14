package rtrb

import (
	"errors"
	"runtime"
	"sync/atomic"
)

// 定义错误类型
var (
	MaxConsumerError  = errors.New("max amount of consumers reached cannot create any more")
	InvalidBufferSize = errors.New("buffer must be of size 2^n")
)

// RingBuffer 是一个支持多消费者的并发安全环形缓冲区
type RingBuffer[T any] struct {
	length            uint32   // 缓冲区总长度（必须是2的幂）
	bitWiseLength     uint32   // 位运算掩码（length-1，用于快速取模）
	headIndex         uint32   // 下一个写入位置索引（原子操作）
	nextReaderIndex   uint32   // 下一个可分配的消费者ID（原子操作）
	maxReaders        int      // 最大消费者数量
	buffer            []T      // 存储数据的环形数组
	readerIndexes     []uint32 // 每个消费者的读取位置索引
	readerActiveFlags []uint32 // 消费者状态标志：0=未使用，1=活跃，2=创建中
}

// CreateBuffer 初始化环形缓冲区
// size 必须是2的幂，maxReaders 指定最大消费者数量
func CreateBuffer[T any](size uint32, maxReaders uint32) (RingBuffer[T], error) {
	if size&(size-1) != 0 { // 检查是否为2的幂
		return RingBuffer[T]{}, InvalidBufferSize
	}

	return RingBuffer[T]{
		buffer:            make([]T, size), // 初始化缓冲区
		length:            size,
		bitWiseLength:     size - 1, // 用于位运算代替取模
		headIndex:         0,
		nextReaderIndex:   0,
		maxReaders:        int(maxReaders),
		readerIndexes:     make([]uint32, maxReaders),
		readerActiveFlags: make([]uint32, maxReaders),
	}, nil
}

// CreateConsumer 创建一个新的消费者
// 通过原子操作分配未使用的消费者槽位，支持并发安全创建
func (buffer *RingBuffer[T]) CreateConsumer() (Consumer[T], error) {
	for readerIndex := range buffer.readerActiveFlags {
		// 尝试将未使用的槽位（0）标记为创建中（2）
		if atomic.CompareAndSwapUint32(&buffer.readerActiveFlags[readerIndex], 0, 2) {
			// 初始化消费者读取位置为当前写入位置
			buffer.readerIndexes[readerIndex] = atomic.LoadUint32(&buffer.headIndex)
			// 标记槽位为活跃状态（1）
			atomic.StoreUint32(&buffer.readerActiveFlags[readerIndex], 1)

			// 更新下一个可用消费者索引（如果当前索引是最后一个）
			atomic.CompareAndSwapUint32(&buffer.nextReaderIndex, uint32(readerIndex), uint32(readerIndex)+1)

			return Consumer[T]{
				id:   uint32(readerIndex),
				ring: buffer,
			}, nil
		}
	}
	return Consumer[T]{}, MaxConsumerError
}

// removeConsumer 移除指定消费者
func (buffer *RingBuffer[T]) removeConsumer(readerId uint32) {
	// 标记槽位为未使用状态
	atomic.StoreUint32(&buffer.readerActiveFlags[readerId], 0)
	// 尝试回退nextReaderIndex（仅当当前移除的是最新消费者时有效）
	atomic.CompareAndSwapUint32(&buffer.nextReaderIndex, readerId+1, readerId)
}

// Write 向缓冲区写入数据（并发安全）
func (buffer *RingBuffer[T]) Write(value T) {
	var offset uint32
	var i uint32

attemptWrite:
	// 获取当前活跃消费者数量
	nextReaderIndex := atomic.LoadUint32(&buffer.nextReaderIndex)

	for i = 0; i < nextReaderIndex; i++ {
		if atomic.LoadUint32(&buffer.readerActiveFlags[i]) == 1 {
			// 计算消费者i的读取位置相对于当前写入位置的偏移
			offset = atomic.LoadUint32(&buffer.readerIndexes[i]) + buffer.length

			// 如果缓冲区已满（某个消费者的读取位置落后一个完整缓冲区长度）
			if offset == buffer.headIndex {
				runtime.Gosched() // 让出CPU时间片
				goto attemptWrite // 重试
			}
		}
	}

	// 计算下一个写入位置并使用位运算取模
	nextIndex := buffer.headIndex + 1
	buffer.buffer[nextIndex&buffer.bitWiseLength] = value
	atomic.StoreUint32(&buffer.headIndex, nextIndex)
}

// readIndex 从指定消费者读取数据
func (buffer *RingBuffer[T]) readIndex(readerIndex uint32) T {
	// 计算下一个要读取的位置
	newIndex := buffer.readerIndexes[readerIndex] + 1

	// 等待数据可用（写入位置已更新）
	for newIndex > atomic.LoadUint32(&buffer.headIndex) {
		runtime.Gosched() // 让出CPU时间片
	}

	// 读取数据并更新消费者读取位置
	value := buffer.buffer[newIndex&buffer.bitWiseLength]
	atomic.StoreUint32(&buffer.readerIndexes[readerIndex], newIndex)
	return value
}

// Consumer 消费者结构体
type Consumer[T any] struct {
	ring *RingBuffer[T] // 关联的环形缓冲区
	id   uint32         // 消费者ID
}

// Remove 销毁当前消费者
func (consumer *Consumer[T]) Remove() {
	consumer.ring.removeConsumer(consumer.id)
}

// Get 从环形缓冲区获取数据（阻塞直到数据可用）
func (consumer *Consumer[T]) Get() T {
	return consumer.ring.readIndex(consumer.id)
}
