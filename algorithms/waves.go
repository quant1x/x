package algorithms

import (
	"fmt"
	"math"
	"slices"
	"sort"
)

const (
	floatEps = 1e-9
	// 其他常量...
)

// ExtremeType 表示极值方向
type ExtremeType int

const (
	ExtremePeak   ExtremeType = iota // 波峰（局部最大值）
	ExtremeTrough                    // 波谷（局部最小值）
)

func (e ExtremeType) String() string {
	switch e {
	case ExtremePeak:
		return "ExtremePeak"
	case ExtremeTrough:
		return "ExtremeTrough"
	default:
		return "Unknown"
	}
}

// SegmentSide 表示自由段的位置（用于 processSegment）
type SegmentSide int

const (
	SideLeft SegmentSide = iota
	SideRight
)

// SearchMode 搜索模式
type SearchMode int

const (
	FindInflection SearchMode = iota // 从左到右：找拐点
	PreserveTrend                    // 从右到左：保终局
)

func (m SearchMode) String() string {
	switch m {
	case FindInflection:
		return "FindInflection"
	case PreserveTrend:
		return "PreserveTrend"
	default:
		return "Unknown"
	}
}

// PeaksResult 返回结果
type PeaksResult struct {
	Peaks     []int // 主趋势波峰（含所有主峰）
	Breakouts []int // 异常突破点
}

func (r PeaksResult) String() string {
	return fmt.Sprintf("Peaks: %v, Breakouts: %v", r.Peaks, r.Breakouts)
}

// SideModes 允许为左侧和右侧自由段独立设置检测模式
type SideModes struct {
	Left  SearchMode // 第一个主峰/主谷左侧使用的模式
	Right SearchMode // 最后一个主峰/主谷右侧使用的模式
}

func processExtremeSegment(
	data []float64,
	segStart, segEnd int,
	extremes []int,
	mode SearchMode,
	results *[]int,
	breakouts *[]int,
	mainVal float64,
	direction ExtremeType,
	side SegmentSide,
) {
	// 增强边界检查：防止越界访问 data[idx]
	if segStart < 0 || segEnd > len(data) || segStart >= segEnd {
		return
	}

	// 收集该段内非主极值的次级极值点
	var segExtremes []int
	for _, idx := range extremes {
		if idx >= segStart && idx < segEnd && math.Abs(data[idx]-mainVal) >= floatEps {
			segExtremes = append(segExtremes, idx)
		}
	}

	if len(segExtremes) == 0 {
		return
	}

	// 根据方向和模式决定遍历顺序与比较逻辑
	var valid []int
	var increasing bool // true: 允许非递减；false: 允许非递增
	var reverseOrder bool

	switch {
	case mode == FindInflection && side == SideLeft && direction == ExtremePeak:
		// 左侧波峰：从右向左，非递增（不能比前一波高）
		reverseOrder = true
		increasing = false
	case mode == FindInflection && side == SideRight && direction == ExtremePeak:
		// 右侧波峰：从左向右，非递增
		reverseOrder = false
		increasing = false
	case mode == PreserveTrend && side == SideLeft && direction == ExtremePeak:
		// 左侧波峰：从左向右，非递减（保持上升趋势）
		reverseOrder = false
		increasing = true
	case mode == PreserveTrend && side == SideRight && direction == ExtremePeak:
		// 右侧波峰：从右向左，非递减 → 反转后时间顺序仍正
		reverseOrder = true
		increasing = true

	// 波谷逻辑（反转比较方向）
	case mode == FindInflection && side == SideLeft && direction == ExtremeTrough:
		// 左侧波谷：从右向左，非递减（不能比前一波低）=> 趋势抬高
		reverseOrder = true
		increasing = true
	case mode == FindInflection && side == SideRight && direction == ExtremeTrough:
		// 右侧波谷：从左向右，非递减
		reverseOrder = false
		increasing = true
	case mode == PreserveTrend && side == SideLeft && direction == ExtremeTrough:
		// 左侧波谷：从左向右，非递增（保持下降趋势）
		reverseOrder = false
		increasing = false
	case mode == PreserveTrend && side == SideRight && direction == ExtremeTrough:
		// 右侧波谷：从右向左，非递增
		reverseOrder = true
		increasing = false
	}

	// 准备遍历顺序
	indices := segExtremes
	if reverseOrder {
		// 逆序遍历
		for i := len(segExtremes) - 1; i >= 0; i-- {
			checkAndAppend(data, segExtremes[i], &valid, breakouts, increasing)
		}
	} else {
		// 正序遍历
		for _, idx := range indices {
			checkAndAppend(data, idx, &valid, breakouts, increasing)
		}
	}

	// 如果是逆序处理，结果需要反转以恢复时间顺序
	if reverseOrder {
		slices.Reverse(valid)
	}

	*results = append(*results, valid...)
}

// checkAndAppend 判断当前点是否符合趋势（非递增/非递减）
func checkAndAppend(
	data []float64,
	currIdx int,
	valid *[]int,
	breakouts *[]int,
	shouldIncrease bool, // true: 非递减；false: 非递增
) {
	if len(*valid) == 0 {
		*valid = append(*valid, currIdx)
		return
	}

	lastVal := data[(*valid)[len(*valid)-1]]
	currVal := data[currIdx]

	if shouldIncrease {
		// 要求非递减：curr >= last
		if currVal >= lastVal {
			*valid = append(*valid, currIdx)
		} else {
			*breakouts = append(*breakouts, currIdx)
		}
	} else {
		// 要求非递增：curr <= last
		if currVal <= lastVal {
			*valid = append(*valid, currIdx)
		} else {
			*breakouts = append(*breakouts, currIdx)
		}
	}
}

// FindExtremesWithBreakouts 在 [start, end) 区间分析波峰或波谷
func FindExtremesWithBreakouts(
	data []float64,
	start, end int,
	modes SideModes,
	direction ExtremeType,
) PeaksResult {
	result := PeaksResult{}

	// 防护：nil、越界、长度不足
	if data == nil || len(data) == 0 ||
		start < 0 || end > len(data) || start >= end || len(data) < 3 {
		return result
	}

	// 1. 找所有极值点（局部 extremum）
	var extremes []int

	// 左端点
	if start+1 < end {
		if (direction == ExtremePeak && data[start] > data[start+1]) ||
			(direction == ExtremeTrough && data[start] < data[start+1]) {
			extremes = append(extremes, start)
		}
	}

	// 内部点
	for i := start + 1; i <= end-2; i++ {
		isPeak := data[i-1] < data[i] && data[i] > data[i+1]
		isTrough := data[i-1] > data[i] && data[i] < data[i+1]

		if (direction == ExtremePeak && isPeak) || (direction == ExtremeTrough && isTrough) {
			extremes = append(extremes, i)
		}
	}

	if len(extremes) == 0 {
		return result
	}

	// 2. 找主极值（全局最值）
	var mainVal float64
	if direction == ExtremePeak {
		mainVal = data[extremes[0]]
		for _, i := range extremes {
			if data[i] > mainVal {
				mainVal = data[i]
			}
		}
	} else {
		mainVal = data[extremes[0]]
		for _, i := range extremes {
			if data[i] < mainVal {
				mainVal = data[i]
			}
		}
	}

	// 3. 收集所有主极值点（使用浮点容差）
	var majorExtremes []int
	for _, i := range extremes {
		if math.Abs(data[i]-mainVal) < floatEps {
			majorExtremes = append(majorExtremes, i)
		}
	}

	// 安全兜底：主极值点为空（理论上不会，但防浮点误差或极端情况）
	if len(majorExtremes) == 0 {
		return result
	}

	sort.Ints(majorExtremes)
	firstMajor := majorExtremes[0]
	lastMajor := majorExtremes[len(majorExtremes)-1]

	var results []int
	var breakouts []int

	// 4. 处理左侧自由段 [start, firstMajor)
	processExtremeSegment(
		data, start, firstMajor,
		extremes, modes.Left,
		&results, &breakouts,
		mainVal, direction, SideLeft,
	)

	// 5. 处理右侧自由段 [lastMajor+1, end)
	processExtremeSegment(
		data, lastMajor+1, end,
		extremes, modes.Right,
		&results, &breakouts,
		mainVal, direction, SideRight,
	)

	// 6. 加入主极值点
	results = append(results, majorExtremes...)
	sort.Ints(results)
	sort.Ints(breakouts)

	result.Peaks = results
	result.Breakouts = breakouts
	return result
}
