package core

import "log/slog"

var (
	logger = slog.Default()
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
