package algorithms

import (
	"fmt"
	"testing"
)

func TestFindPeaksValleys(t *testing.T) {
	highList := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	lowList := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3}
	fmt.Println(len(highList), len(lowList))
	fmt.Println("----------")
	peaks, valleys := FindPeaksValleys(highList, lowList)
	fmt.Println(peaks, valleys)
}

// PeaksResult 返回主波峰和异常突破点
type PeaksResult struct {
	Peaks     []int // 主趋势波峰（含主峰及其合规左右）
	Breakouts []int // 打破单调性的强波峰
}

// FindPeaksWithBreakouts 检测主波峰 + 异常突破点
func FindPeaksWithBreakouts(data []float64) PeaksResult {
	if len(data) == 0 {
		return PeaksResult{}
	}

	// 1. 检测所有局部波峰
	var highs []int
	for i := 1; i < len(data)-1; i++ {
		if data[i-1] < data[i] && data[i] > data[i+1] {
			highs = append(highs, i)
		}
	}
	if len(highs) == 0 {
		return PeaksResult{}
	}

	// 2. 找出全局最大值
	maxVal := data[highs[0]]
	for _, i := range highs {
		if data[i] > maxVal {
			maxVal = data[i]
		}
	}

	// 收集主峰（值等于 maxVal 的波峰）
	var majorPeaks []int
	for _, i := range highs {
		if data[i] == maxVal {
			majorPeaks = append(majorPeaks, i)
		}
	}

	var peaks []int
	var breakouts []int
	start := 0

	for _, peakIdx := range majorPeaks {
		// 找 peakIdx 在 highs 中的位置
		pos := -1
		for i := start; i < len(highs); i++ {
			if highs[i] == peakIdx {
				pos = i
				break
			}
		}
		if pos == -1 {
			continue
		}

		// 左侧 [start, pos)：非递减
		var left []int
		for i := start; i < pos; i++ {
			idx := highs[i]
			if len(left) == 0 || data[idx] >= data[left[len(left)-1]] {
				left = append(left, idx)
			} else {
				// 可选：左侧也记录异常？通常不关心
			}
		}

		// 右侧 (pos, next)：非递增，记录违规者
		var right []int
		j := pos + 1
		for j < len(highs) {
			idx := highs[j]
			// 遇到下一个主峰，停止
			if data[idx] == maxVal {
				break
			}

			// 检查是否满足非递增
			if len(right) == 0 || data[idx] <= data[right[len(right)-1]] {
				right = append(right, idx)
			} else {
				// ❌ 违反非递增！但值显著，记录为 breakout
				// 可加过滤：只有值 > 某阈值才报警，当前无阈值，全部记录
				breakouts = append(breakouts, idx)
				// 不加入 right
			}
			j++
		}

		// 合并主序列
		peaks = append(peaks, left...)
		peaks = append(peaks, peakIdx)
		peaks = append(peaks, right...)

		start = j
	}

	return PeaksResult{
		Peaks:     peaks,
		Breakouts: breakouts,
	}
}

func dataFromIndices(data []float64, indices []int) []float64 {
	var res []float64
	for _, i := range indices {
		res = append(res, data[i])
	}
	return res
}

func TestWavesBasic(t *testing.T) {
	highList := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	lowList := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3}

	fmt.Println(len(highList), len(lowList))
	fmt.Println("----------")

	result := FindPeaksWithBreakouts(highList)

	fmt.Println("主趋势波峰索引:", result.Peaks)
	fmt.Println("主趋势波峰值:  ", dataFromIndices(highList, result.Peaks))

	fmt.Println("异常突破点索引:", result.Breakouts)
	fmt.Println("异常突破点值:  ", dataFromIndices(highList, result.Breakouts))
}
