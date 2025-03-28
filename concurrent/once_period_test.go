package concurrent

import (
	"fmt"
	"github.com/quant1x/x/core"
	"testing"
	"time"
)

func TestPeriodOnce(t *testing.T) {
	once := CreatePeriodOnceWithSecond(1)
	once.Do(func() {
		fmt.Println("1-", time.Now())
	})
	core.WaitForShutdown(5)
	once = CreatePeriodOnceWithHourAndMinute(9, 0)
	once.Do(func() {
		fmt.Println("2-", time.Now())
	})
	core.WaitForShutdown(5)
}
