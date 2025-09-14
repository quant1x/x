package main

import (
	"time"

	logger_error "github.com/quant1x/x/logger-error"
	"go.uber.org/zap"
)

// 主函数示例
func main() {
	//logDir := "./logs"
	//logger.InitLogger(logDir, logger.INFO)
	count := 1000
	//logger.Fatal("This is fatal")
	for i := 0; i < count; i++ {
		// 输出日志
		logger_error.Infof("%d: This is an info message, %+v", i, zap.String("user", "Alice"))
		logger_error.Errorf("%d: This is an error message, %+v", i, zap.Int("code", 500))
		logger_error.Debugf("This is a debug message, %d", i)
		logger_error.Warnf("This is a warn message, %+v", zap.Int("code", 200))
		time.Sleep(1 * time.Second)
	}
}
