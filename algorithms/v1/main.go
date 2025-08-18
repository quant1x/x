package main

import (
	"fmt"
	"sort"
)

// SearchMode 搜索模式
type SearchMode int

const (
	FindInflection SearchMode = iota // 从左到右：找拐点，趋势破坏后停止
	PreserveTrend                    // 从右到左：保终局，只取最右
)

// PeaksResult 返回结果
type PeaksResult struct {
	Peaks     []int // 主趋势波峰（合规）
	Breakouts []int // 打破趋势的异常点
}

// FindPeaksWithBreakouts 在 [start, end) 区间分析波峰
func FindPeaksWithBreakouts(
	data []float64,
	start, end int,
	mode SearchMode,
) PeaksResult {
	result := PeaksResult{}

	// 边界检查
	if start < 0 || end > len(data) || start >= end || len(data) < 3 {
		return result
	}

	// 1. 检测波峰（含左端点）
	var highs []int

	// a. 左端点
	if start+1 < end && data[start] > data[start+1] {
		highs = append(highs, start)
	}

	// b. 内部点
	for i := start + 1; i <= end-2; i++ {
		if data[i-1] < data[i] && data[i] > data[i+1] {
			highs = append(highs, i)
		}
	}

	if len(highs) == 0 {
		return result
	}

	// 2. 找主峰（最大值）
	maxVal := data[highs[0]]
	for _, i := range highs {
		if data[i] > maxVal {
			maxVal = data[i]
		}
	}

	var majorPeaks []int
	for _, i := range highs {
		if data[i] == maxVal {
			majorPeaks = append(majorPeaks, i)
		}
	}
	sort.Ints(majorPeaks) // 从左到右

	var peaks []int
	var breakouts []int
	highsIdx := 0

	// 3. 按主峰顺序处理
	for i := 0; i < len(majorPeaks); i++ {
		peakOrig := majorPeaks[i]

		// 找到 peakOrig 在 highs 中的位置
		pos := -1
		for j := highsIdx; j < len(highs); j++ {
			if highs[j] == peakOrig {
				pos = j
				break
			}
		}
		if pos == -1 {
			continue
		}

		// 左侧：[highsIdx, pos) -> 非递减
		var left []int
		for j := highsIdx; j < pos; j++ {
			idx := highs[j]
			if len(left) == 0 || data[idx] >= data[left[len(left)-1]] {
				left = append(left, idx)
			}
		}

		// 右侧过滤
		rightResult := filterRightPeaks(highs, pos, end, data, maxVal, mode)

		// 合并主趋势
		peaks = append(peaks, left...)
		peaks = append(peaks, peakOrig)
		peaks = append(peaks, rightResult.Peaks...)

		// 收集异常
		breakouts = append(breakouts, rightResult.Breakouts...)

		// 如果是 FindInflection 且趋势破坏，后续主峰不再加入 peaks
		if mode == FindInflection && rightResult.TrendBroken {
			// 将后续所有主峰加入 breakouts
			for k := i + 1; k < len(majorPeaks); k++ {
				breakouts = append(breakouts, majorPeaks[k])
			}
			break // 跳出主峰循环
		}

		// 更新 highsIdx
		if len(rightResult.Peaks) > 0 {
			last := rightResult.Peaks[len(rightResult.Peaks)-1]
			for j := pos + 1; j < len(highs); j++ {
				if highs[j] == last {
					highsIdx = j + 1
					break
				}
			}
		} else {
			highsIdx = pos + 1
		}
	}

	// 排序输出
	sort.Ints(peaks)
	sort.Ints(breakouts)

	result.Peaks = peaks
	result.Breakouts = breakouts
	return result
}

// RightFilterResult 右侧过滤结果
type RightFilterResult struct {
	Peaks       []int
	Breakouts   []int
	TrendBroken bool // 趋势是否破坏
}

// filterRightPeaks 过滤右侧波峰
func filterRightPeaks(
	highs []int,
	pos int,
	end int,
	data []float64,
	maxVal float64,
	mode SearchMode,
) RightFilterResult {
	var result RightFilterResult

	switch mode {
	case FindInflection:
		// 从左到右：非递增，一旦打破，后续全丢
		for j := pos + 1; j < len(highs); j++ {
			idx := highs[j]
			if idx >= end {
				break
			}
			if data[idx] == maxVal {
				break
			}
			if len(result.Peaks) == 0 || data[idx] <= data[result.Peaks[len(result.Peaks)-1]] {
				result.Peaks = append(result.Peaks, idx)
			} else {
				// 趋势破坏
				result.Breakouts = append(result.Breakouts, idx)
				result.TrendBroken = true
				break
			}
		}

	case PreserveTrend:
		// 从右到左：只保留最右一个有效点
		for j := len(highs) - 1; j > pos; j-- {
			idx := highs[j]
			if idx >= end {
				continue
			}
			if data[idx] == maxVal {
				result.Peaks = append(result.Peaks, idx)
				return result
			}
			result.Peaks = append(result.Peaks, idx)
			return result
		}
	}

	return result
}

// dataFromIndices 辅助
func dataFromIndices(data []float64, indices []int) []float64 {
	var res []float64
	for _, i := range indices {
		if i >= 0 && i < len(data) {
			res = append(res, data[i])
		}
	}
	return res
}

// ==================== 测试 ====================

func main() {
	data := []float64{10, 1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}

	fmt.Println("数据:", data)
	fmt.Println("索引: [0 1 2 3 4 5 6 7 8 9 10 11 12]")
	fmt.Println()

	// 示例1：FindInflection
	result1 := FindPeaksWithBreakouts(data, 0, len(data), FindInflection)
	fmt.Println("【FindInflection 模式】")
	fmt.Println("主趋势波峰索引:", result1.Peaks)
	fmt.Println("主趋势波峰值:  ", dataFromIndices(data, result1.Peaks))
	fmt.Println("异常突破点索引:", result1.Breakouts)
	fmt.Println()

	// 示例2：PreserveTrend
	result2 := FindPeaksWithBreakouts(data, 0, len(data), PreserveTrend)
	fmt.Println("【PreserveTrend 模式】")
	fmt.Println("主趋势波峰索引:", result2.Peaks)
	fmt.Println("主趋势波峰值:  ", dataFromIndices(data, result2.Peaks))
	fmt.Println("异常突破点索引:", result2.Breakouts)
}
