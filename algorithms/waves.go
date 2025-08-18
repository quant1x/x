package algorithms

import (
	"slices"
)

// find_extremes 查找波峰波谷
func find_extremes(list []float64, start, end int) (extremes []int) {
	length := len(list)
	if length < 3 {
		return nil
	}
	if length > end {
		return nil
	}
	if start < 0 {
		return nil
	}

	// 使用 int8 节省内存（diff 值仅为 -1,0,1）
	diff := make([]int8, end) // n 个元素，最后一个为 0

	// 1. 计算一阶符号差分（n-1 项），最后一项补 0
	for i := start; i < end-1; i++ {
		// 内联 compare
		if list[i+1] > list[i] {
			diff[i] = 1
		} else if list[i+1] < list[i] {
			diff[i] = -1
		} // else 0（默认）

	}
	// 最后一个差分设为 0（无后续）
	diff[end-1] = 0
	diff[end-1] = 0

	// 2. 处理平台（diff == 0 的点），合并 high 和 low 处理
	for i := start; i < end-1; i++ {
		// 处理 high
		if diff[i] == 0 {
			if i == start {
				for j := start + 1; j < end-1; j++ {
					if diff[j] != 0 {
						diff[i] = diff[j]
						break
					}
				}
			} else if i == end-2 {
				diff[i] = diff[i-1]
			} else {
				diff[i] = diff[i+1]
			}
		}
	}

	// 3. 检测波峰波谷：二阶差分
	// 预分配空间（通常极值点不会超过 n/3）
	extremes = make([]int, 0, end/3)
	for i := start; i < end-1; i++ {
		dHigh := int(diff[i+1]) - int(diff[i])

		if dHigh == -2 {
			extremes = append(extremes, i+1)
		}
	}

	return extremes
}

// FindPeaksValleys 查找波峰波谷
func FindPeaksValleys(highList, lowList []float64) (peaks, valleys []int) {
	n := len(highList)
	if n != len(lowList) || n < 3 {
		return nil, nil
	}

	// 使用 int8 节省内存（diff 值仅为 -1,0,1）
	diffHigh := make([]int8, n) // n 个元素，最后一个为 0
	diffLow := make([]int8, n)

	// 1. 计算一阶符号差分（n-1 项），最后一项补 0
	for i := 0; i < n-1; i++ {
		// 内联 compare
		if highList[i+1] > highList[i] {
			diffHigh[i] = 1
		} else if highList[i+1] < highList[i] {
			diffHigh[i] = -1
		} // else 0（默认）

		if lowList[i+1] > lowList[i] {
			diffLow[i] = 1
		} else if lowList[i+1] < lowList[i] {
			diffLow[i] = -1
		}
	}
	// 最后一个差分设为 0（无后续）
	diffHigh[n-1] = 0
	diffLow[n-1] = 0

	// 2. 处理平台（diff == 0 的点），合并 high 和 low 处理
	for i := 0; i < n-1; i++ {
		// 处理 high
		if diffHigh[i] == 0 {
			if i == 0 {
				for j := 1; j < n-1; j++ {
					if diffHigh[j] != 0 {
						diffHigh[i] = diffHigh[j]
						break
					}
				}
			} else if i == n-2 {
				diffHigh[i] = diffHigh[i-1]
			} else {
				diffHigh[i] = diffHigh[i+1]
			}
		}

		// 处理 low
		if diffLow[i] == 0 {
			if i == 0 {
				for j := 1; j < n-1; j++ {
					if diffLow[j] != 0 {
						diffLow[i] = diffLow[j]
						break
					}
				}
			} else if i == n-2 {
				diffLow[i] = diffLow[i-1]
			} else {
				diffLow[i] = diffLow[i+1]
			}
		}
	}

	// 3. 检测波峰波谷：二阶差分
	// 预分配空间（通常极值点不会超过 n/3）
	peaks = make([]int, 0, n/3)
	valleys = make([]int, 0, n/3)

	for i := 0; i < n-1; i++ {
		dHigh := int(diffHigh[i+1]) - int(diffHigh[i])
		dLow := int(diffLow[i+1]) - int(diffLow[i])

		if dHigh == -2 {
			peaks = append(peaks, i+1)
		}
		if dLow == 2 {
			valleys = append(valleys, i+1)
		}
	}

	return peaks, valleys
}

// ExtremePoint 数据极点, X轴为时间类切片的索引, Y轴为具体数值
type ExtremePoint struct {
	X      int     // Y切片的索引
	Y      float64 // 值
	IsPeak bool    // 是否波峰
}

type WaveDirection int

const (
	WaveLeft  WaveDirection = iota // 左边
	WaveRight                      // 右边
)

type ExtremeType uint8

const (
	ExtremeFlat   ExtremeType = 0         // 平台
	ExtremePeak   ExtremeType = 1 << iota // 波峰
	ExtremeValley                         // 波谷
	ExtremeLatest                         // 最近
	ExtremeGlobal                         // 全局
)

func (e ExtremeType) Compare() func(a, b float64) bool {
	// 定义比较函数：波峰用大于，波谷用小于
	var compare func(a, b float64) bool
	if (e & ExtremePeak) == ExtremePeak {
		compare = func(a, b float64) bool { return a >= b }
	} else if (e & ExtremeValley) == ExtremeValley {
		compare = func(a, b float64) bool { return a <= b }
	} else if e == ExtremeFlat {
		compare = func(a, b float64) bool { return a == b }
	} else {
		panic("unknown ExtremeType")
	}
	return compare
}

func (e ExtremeType) IsPeak() bool {
	return e&ExtremePeak == ExtremePeak
}

func (e ExtremeType) IsValley() bool {
	return e&ExtremeValley == ExtremeValley
}

func (e ExtremeType) Has(other ExtremeType) bool {
	return (e & other) == other
}

func (e ExtremeType) IsLatest() bool {
	return (e & ExtremeLatest) == ExtremeLatest
}

func (e ExtremeType) IsGlobal() bool {
	return (e & ExtremeGlobal) == ExtremeGlobal
}

// find_monotonic_extremes 检测单调序列中的极值点（波峰或波谷）
//
//	支持从左到右或从右到左扫描，返回按原始顺序排列的索引
func find_monotonic_extremes(data []ExtremePoint, direction WaveDirection, typ ExtremeType) []ExtremePoint {
	// 空数据处理
	if len(data) == 0 {
		return []ExtremePoint{}
	}
	// 单点数据
	if len(data) == 1 {
		return data
	}

	var extremes []ExtremePoint
	var startIdx, endIdx, step int

	// 设置扫描方向
	if direction == WaveLeft {
		startIdx, endIdx, step = 0, len(data), 1
	} else if direction == WaveRight {
		startIdx, endIdx, step = len(data)-1, -1, -1
	} else {
		panic("direction must be 'left' or 'right'")
	}

	// 初始值
	prevVal := data[startIdx]

	compare := typ.Compare()

	// 遍历数据
	for currentIdx := startIdx + step; currentIdx != endIdx; currentIdx += step {
		currentVal := data[currentIdx]

		if compare(currentVal.Y, prevVal.Y) {
			// 当前值更极值，更新
			prevVal = currentVal
		} else if len(extremes) > 0 && (prevVal.Y == extremes[len(extremes)-1].Y && prevVal.X == extremes[len(extremes)-1].X) {
			// 跳过已记录的相同极值
			continue
		} else {
			// 当前值不再增长/减少，记录之前的极值
			extremes = append(extremes, prevVal)
		}
	}

	// 处理最后一个极值段
	if len(extremes) == 0 || (compare(prevVal.Y, extremes[len(extremes)-1].Y) && prevVal.X != extremes[len(extremes)-1].X) {
		extremes = append(extremes, prevVal)
	}

	// 如果是右向扫描，反转结果以保持原始顺序
	if direction == WaveRight {
		// 反转切片
		slices.Reverse(extremes)
	}

	return extremes
}

type PhasePoint struct {
	X      int
	Y      int
	IsPeak bool
}

func detect_extremes(raw []float64, extremesType ExtremeType, offset ...int) (phases []PhasePoint, extremes []ExtremePoint) {
	rawCount := len(raw)
	//data := make([]ExtremePoint, dataCount)
	// 确定高点的最大值
	var extremum float64
	if extremesType == ExtremePeak {
		extremum = slices.Max(raw)
	} else {
		extremum = slices.Min(raw)
	}
	offsetStart := 0
	offsetEnd := rawCount
	if len(offset) > 0 {
		offsetStart = offset[0]
	}
	if len(offset) > 1 {
		offsetEnd = offset[1]
	}
	list := find_extremes(raw, offsetStart, offsetEnd)
	extremeCount := len(list)
	if extremeCount == 0 {
		return
	}

	isPeak := extremesType.IsPeak()

	// 按照全局的最大值进行数据分段
	for i := 0; i < extremeCount; i++ {
		x := list[i]
		data := PhasePoint{X: i, Y: x, IsPeak: isPeak}
		if len(phases) == 0 || i+1 == offsetEnd || raw[x] == extremum {
			phases = append(phases, data)
		}
	}
	var tempExtremes []ExtremePoint
	compare := extremesType.Compare()
	for i := 1; i < len(phases); i++ {
		if phases[i-1].Y == phases[i].Y {
			phase := phases[i-1]
			tempExtremes = append(tempExtremes, ExtremePoint{X: phase.Y, Y: raw[phase.Y], IsPeak: phase.IsPeak})
		} else {
			start := phases[i-1].X
			end := phases[i].X
			var direction WaveDirection
			if !compare(raw[phases[i-1].Y], raw[phases[i].Y]) {
				direction = WaveLeft
			} else {
				direction = WaveRight
			}

			if i+1 == len(phases) {
				if end+1 == offsetEnd {
					end += 1
				}
			}
			var tmpList []ExtremePoint
			for j := start; j < end; j++ {
				x := list[j]
				tmpList = append(tmpList, ExtremePoint{X: x, Y: raw[x], IsPeak: isPeak})
			}
			if len(tmpList) < 2 {
				continue
			}
			subExtremes := find_monotonic_extremes(tmpList, direction, extremesType)
			tempExtremes = append(tempExtremes, subExtremes...)
		}
	}

	extremes = append(extremes, tempExtremes...)
	return
}
