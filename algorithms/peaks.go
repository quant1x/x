// peaks.go
package algorithms

import (
	"slices"
	"sort"
)

// SearchMode 搜索模式
type SearchMode int

const (
	FindInflection SearchMode = iota // 从左到右：找拐点
	PreserveTrend                    // 从右到左：保终局
)

// SegmentSide 表示自由段的位置（用于 processSegment）
type SegmentSide int

const (
	SideLeft SegmentSide = iota
	SideRight
)

// PeaksResult 返回结果
type PeaksResult struct {
	Peaks     []int // 主趋势波峰（含所有主峰）
	Breakouts []int // 异常突破点
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

	// 1. 检测所有波峰（含左端点）
	var highs []int

	// 左端点
	if start+1 < end && data[start] > data[start+1] {
		highs = append(highs, start)
	}

	// 内部点
	for i := start + 1; i <= end-2; i++ {
		if data[i-1] < data[i] && data[i] > data[i+1] {
			highs = append(highs, i)
		}
	}

	if len(highs) == 0 {
		// 无波峰，返回
		return result
	}

	// 2. 找出全局最大值
	maxVal := data[highs[0]]
	for _, i := range highs {
		if data[i] > maxVal {
			maxVal = data[i]
		}
	}

	// 3. 收集所有主峰（值等于 maxVal 的波峰）
	var majorPeaks []int
	for _, i := range highs {
		if data[i] == maxVal {
			majorPeaks = append(majorPeaks, i)
		}
	}
	sort.Ints(majorPeaks) // 从左到右

	// ✅ 所有主峰必须进入 peaks
	var peaks []int
	var breakouts []int

	// 4. 分段处理：只处理非主峰之间的“自由段”
	// 自由段1: [start, majorPeaks[0]) —— 第一个主峰左侧
	processSegment(data, start, majorPeaks[0], highs, mode, &peaks, &breakouts, maxVal, SideLeft)

	// 自由段2: [lastMajor+1, end) —— 最后一个主峰右侧
	lastMajor := majorPeaks[len(majorPeaks)-1]
	processSegment(data, lastMajor+1, end, highs, mode, &peaks, &breakouts, maxVal, SideRight)

	// 5. 将所有主峰加入 peaks
	peaks = append(peaks, majorPeaks...)

	// 6. 排序输出
	sort.Ints(peaks)
	sort.Ints(breakouts)

	result.Peaks = peaks
	result.Breakouts = breakouts
	return result
}

func processSegment(
	data []float64,
	segStart, segEnd int,
	highs []int,
	mode SearchMode,
	peaks *[]int,
	breakouts *[]int,
	maxVal float64,
	side SegmentSide, // 改为枚举类型
) {
	if segStart >= segEnd {
		return
	}

	// 收集该段内的次级波峰（非主峰）
	var segHighs []int
	for _, idx := range highs {
		if idx >= segStart && idx < segEnd && data[idx] != maxVal {
			segHighs = append(segHighs, idx)
		}
	}

	if len(segHighs) == 0 {
		return
	}

	switch mode {
	case FindInflection:
		if side == SideLeft {
			// 左侧：从主峰往左看，右→左，非递增（一波不比一波高）
			var valid []int
			for i := len(segHighs) - 1; i >= 0; i-- {
				idx := segHighs[i]
				if len(valid) == 0 || data[idx] <= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
					//break // 趋势破坏
				}
			}
			// 反转，使为时间顺序
			slices.Reverse(valid)
			*peaks = append(*peaks, valid...)
		} else {
			// 右侧：从主峰往右看，左→右，非递增（一波不比一波高）
			var valid []int
			for _, idx := range segHighs {
				if len(valid) == 0 || data[idx] <= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
					//break
				}
			}
			*peaks = append(*peaks, valid...)
		}

	case PreserveTrend:
		if side == SideLeft {
			// 左侧：从低到高，左→右，非递减（一波不比一波低）
			var valid []int
			for _, idx := range segHighs {
				if len(valid) == 0 || data[idx] >= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
					//break
				}
			}
			*peaks = append(*peaks, valid...)
		} else {
			// 右侧：从低到高，右→左，非递减（一波不比一波低）
			var valid []int
			for i := len(segHighs) - 1; i >= 0; i-- {
				idx := segHighs[i]
				if len(valid) == 0 || data[idx] >= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
					//break
				}
			}
			// 反转，使为时间顺序
			slices.Reverse(valid)
			*peaks = append(*peaks, valid...)
		}
	}
}
