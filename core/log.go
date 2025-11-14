package core

// Logger 日志接口，避免循环引用
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

var globalLogger Logger = nil

// SetLogger 设置全局logger
func SetLogger(logger Logger) {
	globalLogger = logger
}