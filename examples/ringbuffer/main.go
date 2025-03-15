package main

import (
	"fmt"
	"github.com/quant1x/x/ringbuffer"
	"log"
	"slices"
	"sync"
	"sync/atomic"
)

func main() {
	rb, err := ringbuffer.New[int](1024)
	if err != nil {
		log.Fatal(err)
	}
	var data []int
	var result []int
	var m sync.Mutex
	dataTotal := 1000
	producterNum := 3
	consumerNum := 1
	wgLocal := sync.WaitGroup{}

	prodAppend := func(waitGroup *sync.WaitGroup, v int) {
		m.Lock()
		defer m.Unlock()
		defer waitGroup.Done()
		data = append(data, v)
	}
	conAppend := func(waitGroup *sync.WaitGroup, v int) {
		m.Lock()
		defer m.Unlock()
		defer waitGroup.Done()
		result = append(result, v)
	}

	needReadCount := producterNum * dataTotal
	var readNum atomic.Int32
	readNum.Store(0)
	// 启动4个消费者
	for i := 0; i < consumerNum; i++ {
		wgLocal.Add(1)
		go func(waitGroup *sync.WaitGroup, no int) {
			defer waitGroup.Done()
			//defer c.Close()
			for readNum.Load() < int32(needReadCount) {
				v, err := rb.Read()
				//fmt.Println("consumer:", no, "<=", v)
				if err == nil {
					waitGroup.Add(1)
					go conAppend(waitGroup, v)
					readNum.Add(1)
				} else {
					fmt.Println(err)
					break
				}
				if readNum.Load() > int32(needReadCount)-2 {
					fmt.Println("No:", i, "readNum:", readNum.Load())
				}
			}
			fmt.Println("No:", i, "exit")
			rb.Close()
		}(&wgLocal, i)
	}

	// 启动4个生产者
	for i := 0; i < producterNum; i++ {
		wgLocal.Add(1)
		go func(waitGroup *sync.WaitGroup, no int) {
			defer waitGroup.Done()
			for j := 0; j < dataTotal; j++ {
				v := no*dataTotal + j
				rb.Write(v) // 忽略错误处理
				waitGroup.Add(1)
				go prodAppend(waitGroup, v)
			}
		}(&wgLocal, i)
	}
	wgLocal.Wait()
	slices.Sort(data)
	slices.Sort(result)
	//got := reflect.DeepEqual(data, result)
	//if got != true {
	//	t.Fatalf("want %v, got %v", true, got)
	//}
	fmt.Println(data)
	fmt.Println(result)
	if len(data) != len(result) {
		log.Fatal("len(data) != len(result):", len(data), len(result))
	}
	for i, _ := range data {
		if data[i] != result[i] {
			log.Fatal("data[i] != result[i]:", data[i], result[i])
		}
	}

}
