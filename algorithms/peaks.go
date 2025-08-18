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

// SideModes 允许为左侧和右侧自由段独立设置检测模式
type SideModes struct {
	Left  SearchMode // 第一个主峰/主谷左侧使用的模式
	Right SearchMode // 最后一个主峰/主谷右侧使用的模式
}

// FindPeaksWithBreakouts 在 [start, end) 区间分析波峰
// FindPeaksWithBreakouts 支持左右段独立模式
func FindPeaksWithBreakouts(
	data []float64,
	start, end int,
	modes SideModes,
) PeaksResult {
	result := PeaksResult{}

	if start < 0 || end > len(data) || start >= end || len(data) < 3 {
		return result
	}

	// 1. 找波峰
	var highs []int
	if start+1 < end && data[start] > data[start+1] {
		highs = append(highs, start)
	}
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
	sort.Ints(majorPeaks)

	firstMajor := majorPeaks[0]
	lastMajor := majorPeaks[len(majorPeaks)-1]

	var peaks []int
	var breakouts []int

	// ✅ 左侧段：使用用户指定的 modes.Left
	processSegment(
		data,
		start,
		firstMajor,
		highs,
		modes.Left, // ← 用户自由设置的模式
		&peaks,
		&breakouts,
		maxVal,
		SideLeft,
	)

	// ✅ 右侧段：使用用户指定的 modes.Right
	processSegment(
		data,
		lastMajor+1,
		end,
		highs,
		modes.Right, // ← 用户自由设置的模式
		&peaks,
		&breakouts,
		maxVal,
		SideRight,
	)

	// 加入主峰
	peaks = append(peaks, majorPeaks...)
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

// FindValleysWithBreakouts 在 [start, end) 区间分析波谷
func FindValleysWithBreakouts(
	data []float64,
	start, end int,
	modes SideModes,
) PeaksResult {
	result := PeaksResult{}

	// 边界检查
	if start < 0 || end > len(data) || start >= end || len(data) < 3 {
		return result
	}

	// 1. 检测所有波谷（含左端点）
	var lows []int

	// 左端点
	if start+1 < end && data[start] < data[start+1] {
		lows = append(lows, start)
	}

	// 内部点：严格局部最小值
	for i := start + 1; i <= end-2; i++ {
		if data[i-1] > data[i] && data[i] < data[i+1] {
			lows = append(lows, i)
		}
	}

	if len(lows) == 0 {
		return result
	}

	// 2. 找出全局最小值（波谷中的最低点）
	minVal := data[lows[0]]
	for _, i := range lows {
		if data[i] < minVal {
			minVal = data[i]
		}
	}

	// 3. 收集所有主谷（值等于 minVal 的波谷）
	var majorValleys []int
	for _, i := range lows {
		if data[i] == minVal {
			majorValleys = append(majorValleys, i)
		}
	}
	sort.Ints(majorValleys) // 从左到右排序

	// ✅ 所有主谷必须进入 valleys
	var valleys []int
	var breakouts []int

	// 4. 分段处理自由段：仅处理非主谷之间的区域
	// 自由段1: [start, majorValleys[0]) —— 第一个主谷左侧
	processValleySegment(data, start, majorValleys[0], lows, modes.Left, &valleys, &breakouts, minVal, SideLeft)

	// 自由段2: [lastMajor+1, end) —— 最后一个主谷右侧
	lastMajor := majorValleys[len(majorValleys)-1]
	processValleySegment(data, lastMajor+1, end, lows, modes.Right, &valleys, &breakouts, minVal, SideRight)

	// 5. 将所有主谷加入结果
	valleys = append(valleys, majorValleys...)

	// 6. 排序输出
	sort.Ints(valleys)
	sort.Ints(breakouts)

	result.Peaks = valleys // 复用 PeaksResult，但实际是 valleys
	result.Breakouts = breakouts
	return result
}

func processValleySegment(
	data []float64,
	segStart, segEnd int,
	lows []int,
	mode SearchMode,
	valleys *[]int,
	breakouts *[]int,
	minVal float64,
	side SegmentSide,
) {
	if segStart >= segEnd {
		return
	}

	// 收集该段内的次级波谷（非主谷）
	var segLows []int
	for _, idx := range lows {
		if idx >= segStart && idx < segEnd && data[idx] != minVal {
			segLows = append(segLows, idx)
		}
	}

	if len(segLows) == 0 {
		return
	}

	switch mode {
	case FindInflection:
		if side == SideLeft {
			// 左侧：从主谷往左看，右→左，非递减（一波不比一波低）=> 趋势是“抬高”的反转
			var valid []int
			for i := len(segLows) - 1; i >= 0; i-- {
				idx := segLows[i]
				if len(valid) == 0 || data[idx] >= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
				}
			}
			// 反转为时间顺序（从左到右）
			slices.Reverse(valid)
			*valleys = append(*valleys, valid...)
		} else {
			// 右侧：从主谷往右看，左→右，非递减（一波不比一波低）
			var valid []int
			for _, idx := range segLows {
				if len(valid) == 0 || data[idx] >= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
				}
			}
			*valleys = append(*valleys, valid...)
		}

	case PreserveTrend:
		if side == SideLeft {
			// 左侧：从高到低，左→右，非递增（一波不比一波高）=> 趋势是“下降”的
			var valid []int
			for _, idx := range segLows {
				if len(valid) == 0 || data[idx] <= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
				}
			}
			*valleys = append(*valleys, valid...)
		} else {
			// 右侧：从高到低，右→左，非递增（一波不比一波高）
			var valid []int
			for i := len(segLows) - 1; i >= 0; i-- {
				idx := segLows[i]
				if len(valid) == 0 || data[idx] <= data[valid[len(valid)-1]] {
					valid = append(valid, idx)
				} else {
					*breakouts = append(*breakouts, idx)
				}
			}
			// 反转为时间顺序
			slices.Reverse(valid)
			*valleys = append(*valleys, valid...)
		}
	}
}
