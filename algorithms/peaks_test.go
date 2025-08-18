// peaks_test.go
package algorithms

import (
	"fmt"
	"testing"
)

func TestWavesBasic(t *testing.T) {
	data := []float64{2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}

	t.Log("数据:", data)
	t.Log("索引: [0 1 2 3 4 5 6 7 8 9 10 11 12 13]")
	t.Log("")

	// 示例1：FindInflection
	result1 := FindPeaksWithBreakouts(data, 0, len(data), FindInflection)
	t.Log("【FindInflection 模式】")
	t.Log("主趋势波峰索引:", result1.Peaks)
	t.Log("主趋势波峰值:  ", dataFromIndices(data, result1.Peaks))
	t.Log("异常突破点索引:", result1.Breakouts)
	t.Log("")

	// 示例2：PreserveTrend
	result2 := FindPeaksWithBreakouts(data, 0, len(data), PreserveTrend)
	t.Log("【PreserveTrend 模式】")
	t.Log("主趋势波峰索引:", result2.Peaks)
	t.Log("主趋势波峰值:  ", dataFromIndices(data, result2.Peaks))
	t.Log("异常突破点索引:", result2.Breakouts)

	// 断言
	expected1 := []int{1, 3, 5}
	if !equal(result1.Peaks, expected1) {
		t.Errorf("FindInflection: 期望 %v, 实际 %v", expected1, result1.Peaks)
	}

	expected2 := []int{1, 7, 9}
	if !equal(result2.Peaks, expected2) {
		t.Errorf("PreserveTrend: 期望 %v, 实际 %v", expected2, result2.Peaks)
	}

	lows := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3}
	result21 := FindValleysWithBreakouts(lows, 0, len(lows), FindInflection)
	fmt.Println(result21.Peaks)
	result22 := FindValleysWithBreakouts(lows, 0, len(lows), PreserveTrend)
	fmt.Println(result22.Peaks)
}

func dataFromIndices(data []float64, indices []int) []float64 {
	var res []float64
	for _, i := range indices {
		if i >= 0 && i < len(data) {
			res = append(res, data[i])
		}
	}
	return res
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
