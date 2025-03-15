package concurrent

import (
	"sync"
	"sync/atomic"
	"time"
)

// 状态机标志位 (原子操作)
const (
	stateNeedInit uint32 = 0 // 需要初始化
	stateRunning  uint32 = 1 // 初始化中
	stateDone     uint32 = 2 // 初始化完成
)

type DailyOnce struct {
	state    uint32       // 原子状态标志
	nextTime atomic.Int64 // 下次重置时间(unix nano)
	mu       sync.Mutex   // 仅用于初始化互斥
}

func NewDailyOnce() *DailyOnce {
	d := &DailyOnce{}
	d.nextTime.Store(0) // 初始为需要初始化状态
	return d
}

// 获取数据（自动处理每日重置）
func (d *DailyOnce) Do(initFunc func()) {
	// 快速路径检查
	if atomic.LoadUint32(&d.state) == stateDone {
		if now := time.Now().UnixNano(); now < d.nextTime.Load() {
			return
		}
	}

	// 进入慢速路径
	d.doSlow(initFunc)
}

func (d *DailyOnce) doSlow(initFunc func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 双重检查
	now := time.Now()
	if d.state == stateDone && now.UnixNano() < d.nextTime.Load() {
		return
	}

	// 原子状态转换
	atomic.StoreUint32(&d.state, stateRunning)
	defer atomic.StoreUint32(&d.state, stateDone)

	// 执行初始化
	initFunc()

	// 计算并存储下次重置时间
	next := nextExecutionTime(now)
	d.nextTime.Store(next.UnixNano())
}

// 计算下一个有效时间点（立即执行或次日9点）
func nextExecutionTime(now time.Time) time.Time {
	today9 := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if now.After(today9) {
		return today9.Add(24 * time.Hour)
	}
	return today9
}
