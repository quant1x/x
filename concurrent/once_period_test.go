package concurrent

import (
	"fmt"
	"github.com/quant1x/x/core"
	"testing"
	"time"
)

func TestPeriodOnce(t *testing.T) {
	once, err := CreatePeriodOnce(0, 5)
	if err != nil {
		t.Fatal(err)
	}
	once.Do(func() {
		fmt.Println("1-", time.Now())
	})
	core.WaitForShutdown()
}
