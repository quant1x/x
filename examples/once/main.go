package main

import (
	"fmt"
	"time"

	"github.com/quant1x/x/concurrent"
	"github.com/quant1x/x/core"
)

var (
	num  = 0
	once *concurrent.PeriodOnce
)

func getNum() int {
	once.Do(func() {
		num++
	})
	return num
}

func main() {
	once = concurrent.CreatePeriodOnceWithSecond(5)

	go func() {
		for {
			fmt.Printf("demo: %d\n", getNum())
			time.Sleep(time.Second)
		}
	}()
	core.WaitForShutdown()
}
