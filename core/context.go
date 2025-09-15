package core

import (
	"context"
	"sync"
	"time"

	"github.com/quant1x/x/std/signal"
)

var (
	initOnce          sync.Once
	rootContext       context.Context
	rootCancel        context.CancelFunc
	shutdownWaitGroup sync.WaitGroup
)

// initialize 初始化根上下文和取消函数
func initialize() {
	rootContext, rootCancel = context.WithCancel(context.Background())
}

// Context 获取全局顶层context
func Context() context.Context {
	initOnce.Do(initialize)
	return rootContext
}

// Shutdown 关闭应用程序, 通知所有协程退出
func Shutdown() {
	initOnce.Do(initialize)
	if rootCancel != nil {
		rootCancel()
	}
}

// GetContextWithCancel 创建一个新的可取消上下文
func GetContextWithCancel() (context.Context, context.CancelFunc) {
	initOnce.Do(initialize)

	ctx, cancel := context.WithCancel(rootContext)
	shutdownWaitGroup.Add(1)

	// 包装cancel函数以确保WaitGroup计数正确
	var once sync.Once
	wrappedCancel := func() {
		once.Do(func() {
			cancel()
			shutdownWaitGroup.Done()
		})
	}

	return ctx, wrappedCancel
}

// RegisterHook 注册系统退出的hook
func RegisterHook(name string, callback func()) context.Context {
	ctx, cancel := GetContextWithCancel()

	go func() {
		defer cancel() // 确保无论如何都会调用cancel

		<-ctx.Done() // 等待关闭信号
		callback()   // 执行回调
	}()
	_ = name

	return ctx
}

// applicationShutdown 执行应用退出前的清理工作
func applicationShutdown() {
	if rootCancel != nil {
		rootCancel()
	}
	shutdownWaitGroup.Wait()
}

// WaitForShutdown 阻塞等待关闭信号
//
// 参数d为等待的毫秒数：
//   - 如果传入d > 0，等待指定毫秒后关闭
//   - 如果传入d == 0或未传参，等待中断信号
func WaitForShutdown(d ...int) {
	initOnce.Do(initialize)

	interrupt := signal.NotifyForShutdown()
	delay := 0
	if len(d) > 0 {
		delay = d[0]
	}

	if delay != 0 {
		select {
		case <-time.After(time.Millisecond * time.Duration(delay)):
		case <-rootContext.Done():
		case <-interrupt:
		}
	} else {
		select {
		case <-rootContext.Done():
		case sig := <-interrupt:
			_ = sig // 忽略信号值
		}
	}

	applicationShutdown()
}
