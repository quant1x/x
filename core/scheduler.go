package core

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var (
	//global_context, global_cancel = context.WithCancel(context.Background())
	global_scheduler_cron     gocron.Scheduler
	global_scheduler_once     sync.Once
	global_scheduler_location = time.Local
	//scheduler_done            = make(chan struct{})
)

// 初始化全局调度器
func lazyInitScheduler() {
	// 创建一个调度器
	cron, err := gocron.NewScheduler(
		gocron.WithLocation(global_scheduler_location),
		//gocron.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}
	global_scheduler_cron = cron
	global_scheduler_cron.Start()
	_ = RegisterHook("scheduler", stopScheduler)
}

func stopScheduler() {
	//fmt.Println("stop scheduler-1")
	err := global_scheduler_cron.Shutdown()
	if err != nil {
		//logger.Error("x/cron: scheduler shutdown err:", err)
	}
	//fmt.Println("stop scheduler-2")
}

func AddJob(spec string, cmd func()) error {
	global_scheduler_once.Do(lazyInitScheduler)
	jd := gocron.CronJob(spec, true)
	task := gocron.NewTask(cmd)
	// 添加一个job到全局调度器
	job, err := global_scheduler_cron.NewJob(jd, task)
	if err != nil {
		return err
	}
	_ = job
	return nil
}
