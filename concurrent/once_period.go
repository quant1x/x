package concurrent

import (
	"fmt"
	"github.com/quant1x/x/core"
	"sync"
	"sync/atomic"
)

// PeriodOnce 周期性懒加载锁
type PeriodOnce struct {
	done atomic.Uint32
	m    sync.Mutex
}

// CreatePeriodOnceWithHourAndMinute 创建每日按时分初始化的周期性懒加载锁
func CreatePeriodOnceWithHourAndMinute(hour, minute int) *PeriodOnce {
	spec := fmt.Sprintf("0 %d %d * * *", minute, hour)
	return createPeriodOnce(spec)
}

// CreatePeriodOnceWithSecond 创建间隔秒数初始化的周期性懒加载锁
func CreatePeriodOnceWithSecond(seconds int) *PeriodOnce {
	spec := fmt.Sprintf("*/%d * * * * *", seconds)
	return createPeriodOnce(spec)
}

func createPeriodOnce(spec string) *PeriodOnce {
	once := PeriodOnce{}
	err := core.AddJob(spec, func() {
		once.done.Store(stateNeedInit)
	})
	if err != nil {
		panic(err)
	}
	return &once
}

func (o *PeriodOnce) Do(f func()) {
	if o.done.Load() == 0 {
		o.doSlow(f)
	}
}

func (o *PeriodOnce) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done.Load() == 0 {
		defer o.done.Store(1)
		f()
	}
}
