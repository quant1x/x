package main

import (
	"fmt"
	"github.com/quant1x/x/concurrent"
	"github.com/quant1x/x/core"
	"log"
	"time"
)

func main() {
	once, err := concurrent.CreatePeriodOnce(0, 5)
	if err != nil {
		log.Fatal(err)
	}
	once.Do(func() {
		fmt.Println("1-", time.Now())
	})
	core.WaitForShutdown()
}
