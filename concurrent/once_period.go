package concurrent

import (
	"fmt"
	"github.com/quant1x/x/core"
	"sync"
	"sync/atomic"
	"time"
)

type PeriodOnce struct {
	done atomic.Uint32
	m    sync.Mutex
}

func CreatePeriodOnce(hour, minute int) (*PeriodOnce, error) {
	once := &PeriodOnce{}
	spec := fmt.Sprintf("1/%d %d * * * *", minute, hour)
	spec = fmt.Sprintf("*/%d * * * * *", minute)
	err := core.AddJob(spec, func() {
		fmt.Println("2-", time.Now())
		once.done.Store(stateNeedInit)
	})
	return once, err
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
