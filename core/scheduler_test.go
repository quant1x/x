package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func TestRawScheduler(t *testing.T) {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
				time.Sleep(2 * time.Second)
				fmt.Println(time.Now())
			},
			"hello",
			1,
		),
		// 如果作业已在运行，则 WithSingletonMode 可防止作业再次运行。 这对于不应重叠的作业很有用，并且偶尔 （但不一致）运行时间长于作业运行之间的间隔。
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}
