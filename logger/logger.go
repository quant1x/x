package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/quant1x/x/core"
	"github.com/quant1x/x/rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	Level         zapcore.Level // 日志级别
	Path          string        // 路径
	EnableConsole bool          // 控制台开关
	MaxAge        time.Duration // 最大保留时间
	RotationTime  time.Duration // 日志切割时间
	BufferSize    int           // 缓冲区大小, 单位KB
	FlushInterval time.Duration // 定时刷新间隔, 单位秒
}

var (
	// --------------------------------------------
	// 1. 定义纯文本编码器
	// --------------------------------------------
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	textEncoder = zapcore.NewConsoleEncoder(encoderConfig)
)

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	OFF
	FATAL
)

var (
	defaultLevel = DEBUG
	cfg          = Config{
		Level:         zapcore.DebugLevel,
		MaxAge:        7 * 24 * time.Hour,
		RotationTime:  24 * time.Hour,
		BufferSize:    256,
		FlushInterval: 5,
	}
	logger *zap.SugaredLogger = nil
)

var (
	mu              sync.Mutex
	bufferedWriters []*zapcore.BufferedWriteSyncer
)

func addBufferWriteSyncer(bw *zapcore.BufferedWriteSyncer) {
	if bw == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	bufferedWriters = append(bufferedWriters, bw)
}

// SetLevel 在临时路径记录日志
func SetLevel(level LogLevel) {
	InitLogger("", level)
}

// IsDebug 是否debug日志模式
func IsDebug() bool {
	return cfg.Level == zapcore.DebugLevel
}

// InitLogger 初始化全局日志模块
func InitLogger(path string, level LogLevel) {
	path = strings.TrimSpace(path)
	if path == "" {
		path = os.TempDir()
	}
	defaultLevel = level
	cfg.EnableConsole = false
	switch level {
	case DEBUG:
		cfg.Level = zapcore.DebugLevel
		cfg.EnableConsole = true
	case INFO:
		cfg.Level = zapcore.InfoLevel
	case ERROR:
		cfg.Level = zapcore.ErrorLevel
	case WARN:
		cfg.Level = zapcore.WarnLevel
	case OFF:
		cfg.Level = zapcore.FatalLevel
	case FATAL:
		cfg.Level = zapcore.FatalLevel
	default:
		cfg.Level = zapcore.FatalLevel
	}
	cfg.Path = getLogRoot(path)
	fmt.Println(cfg)
	zapLogger := NewTextLoggerWithCompression(cfg)
	logger = zapLogger.Sugar()
}

func getLogRoot(path string) string {
	applicationName := getApplicationName()
	return filepath.Join(path, applicationName)
}

// getApplicationName 获取执行文件名
func getApplicationName() string {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	arr := strings.Split(exec, ".")
	__applicationName := arr[0]
	return __applicationName
}

var (
	mapLevelToFilename = map[zapcore.Level]string{
		zapcore.DebugLevel:  "debug",
		zapcore.InfoLevel:   "info",
		zapcore.WarnLevel:   "warn",
		zapcore.ErrorLevel:  "error",
		zapcore.DPanicLevel: "fatal",
		zapcore.PanicLevel:  "fatal",
		zapcore.FatalLevel:  "fatal",
	}
	console            = zapcore.AddSync(os.Stdout)
	loggerShutdownOnce sync.Once
)

func getLogger(cfg Config, level zapcore.Level) (zapcore.Core, error) {
	filename, ok := mapLevelToFilename[level]
	if !ok {
		panic("invalid log level")
	}
	// 配置日志滚动器，按天切割
	path := filepath.Join(cfg.Path, filename+"_%Y%m%d.log")
	rl, err := rotatelogs.New(
		path,                              // 文件名格式，带日期
		rotatelogs.WithMaxAge(cfg.MaxAge), // 保留7天的日志
		rotatelogs.WithRotationTime(cfg.RotationTime), // 每24小时切割一次
		rotatelogs.WithHandler(rotatelogs.HandlerFunc(
			func(e rotatelogs.Event) {
				if e.Type() == rotatelogs.FileRotatedEventType {
					if fre, ok := e.(*rotatelogs.FileRotatedEvent); ok {
						oldFilename := fre.PreviousFile()
						if oldFilename == "" {
							return
						}
						compressOldLogs(oldFilename)
					}
				}
			})),
	)

	if err != nil {
		return nil, err
	}
	writeSyncer := zapcore.AddSync(rl)
	// 带缓冲的 WriteSyncer（缓冲区大小 256KB）
	bufferedWriteSyncer := &zapcore.BufferedWriteSyncer{
		WS:            writeSyncer,
		Size:          cfg.BufferSize * 1024,           // 缓冲区大小
		FlushInterval: cfg.FlushInterval * time.Second, // 定时刷新间隔
	}
	var syncers []zapcore.WriteSyncer
	syncers = append(syncers, bufferedWriteSyncer)
	addBufferWriteSyncer(bufferedWriteSyncer)
	if cfg.EnableConsole {
		syncers = append(syncers, console)
	}
	core := zapcore.NewCore(
		textEncoder,
		zapcore.NewMultiWriteSyncer(syncers...),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == level
		}),
	)
	return core, nil
}

// 压缩旧日志文件的钩子函数
func compressOldLogs(previousFile string) {
	const logExt = ".log"
	if filepath.Ext(previousFile) == logExt {
		src, err := os.Open(previousFile)
		if err != nil {
			return
		}
		// 压缩文件：原文件 → 原文件.gz
		gzPath := previousFile[:len(previousFile)-len(logExt)] + ".gz"
		dst, err := os.Create(gzPath)
		if err != nil {
			src.Close()
			return
		}

		gzWriter := gzip.NewWriter(dst)
		fileStat, err := src.Stat()
		if err != nil {
			src.Close()
			dst.Close()
			gzWriter.Close()
			return
		}

		gzWriter.Name = fileStat.Name()
		gzWriter.ModTime = fileStat.ModTime()
		_, err = io.Copy(gzWriter, src) // 压缩内容
		if err != nil {
			src.Close()
			dst.Close()
			gzWriter.Close()
			return
		}

		_ = src.Close()
		_ = dst.Close()
		_ = gzWriter.Close()
		_ = os.Remove(previousFile) // 删除原文件
	}
}

// NewTextLoggerWithCompression 初始化支持压缩的纯文本日志配置
func NewTextLoggerWithCompression(cfg Config) *zap.Logger {
	var cores []zapcore.Core

	if cfg.Level <= zapcore.DebugLevel {
		debugLogger, err := getLogger(cfg, zap.DebugLevel)
		if err != nil {
			panic(err)
		}
		cores = append(cores, debugLogger)
	}

	if cfg.Level <= zapcore.InfoLevel {
		infoLogger, err := getLogger(cfg, zap.InfoLevel)
		if err != nil {
			panic(err)
		}
		cores = append(cores, infoLogger)
	}

	if cfg.Level <= zapcore.WarnLevel {
		warnLogger, err := getLogger(cfg, zap.WarnLevel)
		if err != nil {
			panic(err)
		}
		cores = append(cores, warnLogger)
	}

	if cfg.Level <= zapcore.ErrorLevel {
		errorLogger, err := getLogger(cfg, zap.ErrorLevel)
		if err != nil {
			panic(err)
		}
		cores = append(cores, errorLogger)
	}

	if cfg.Level <= zapcore.FatalLevel {
		fatalLogger, err := getLogger(cfg, zap.FatalLevel)
		if err != nil {
			panic(err)
		}
		cores = append(cores, fatalLogger)
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

//// =========================
//// 🚨 修复点1：定义 Errorf 和 Infof（因为 logger 是 SugaredLogger，不能直接用 Errorf）
//// =========================
//
//func Errorf(format string, args ...interface{}) {
//	if logger != nil {
//		logger.Errorf(format, args...)
//	} else {
//		fmt.Printf("[ERROR] "+format+"\n", args...)
//	}
//}
//
//func Infof(format string, args ...interface{}) {
//	if logger != nil {
//		logger.Infof(format, args...)
//	} else {
//		fmt.Printf("[INFO] "+format+"\n", args...)
//	}
//}

func init() {
	tempPath := os.TempDir()
	InitLogger(tempPath, defaultLevel)

	//chSignal := signal.NotifyForShutdown()
	//go waitForStop(chSignal)
	core.RegisterHook("logger", Shutdown)
}

func loggerSyncAndShutdown() {
	// 1. 停止所有自定义缓冲写入器（如 BufferedWriteSyncer）
	for i, bw := range bufferedWriters {
		if err := bw.Stop(); err != nil {
			Errorf("业务写入器 [%d] 停止失败: %v", i, err)
		}
	}

	// 2. 优雅刷新日志 —— 必须检查错误！
	if logger != nil {
		l := logger.Desugar()            // 👈 关键：获取底层 *zap.Logger
		if err := l.Sync(); err != nil { // 👈 正确方式！
			Errorf("日志同步失败，可能导致部分日志丢失: %v", err)
		} else {
			Infof("日志已成功刷新并关闭")
		}
	}
}

// Shutdown 优雅关闭日志系统，刷盘并释放资源
// 可被 main、测试、HTTP 接口等主动调用
func Shutdown() {
	loggerShutdownOnce.Do(loggerSyncAndShutdown)
}
