package mpb

import (
	"fmt"
	"strings"
	"sync"
)

var (
	mu           sync.Mutex
	gSrcLine     = 0 //起点行
	gCurrentLine = 0 //当前行
	gMaxLine     = 0 //最大行
)

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	gSrcLine = 0
	gCurrentLine = 0
	gMaxLine = 0
}

func GetMaxLine() int {
	return gMaxLine
}

func SetMaxLine(line int) {
	mu.Lock()
	defer mu.Unlock()
	gMaxLine = line
}

// 调整最大行数, 补充新的空白行给新进度条
func adjustLine(line int) int {
	mu.Lock()
	defer mu.Unlock()
	old := gMaxLine
	if line <= 0 {
		gMaxLine++
		line = gMaxLine
	}
	if line > gMaxLine {
		gMaxLine = line
	}
	if old > 0 && gMaxLine > old {
		fmt.Printf(strings.Repeat("\r\n", gMaxLine-old))
	}
	_ = old
	return line
}

// 移动光标到指定的进度条的行号
func barMove(line int) {
	fmt.Printf("\033[%dA\033[%dB", gCurrentLine, line)
	gCurrentLine = line
}

// 更新指定行的进度条信息
func barPrintf(line int, format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()

	barMove(line)
	fmt.Printf("\r"+format, args...)
	barMove(gMaxLine)
}
