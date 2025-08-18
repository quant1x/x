// peaks_test.go
package algorithms

import (
	"fmt"
	"strings"
	"testing"
)

func TestWavesV2(t *testing.T) {
	data := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	lows := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3} // 注意：len(lows) == 13

	// 辅助函数：索引对齐输出
	printHeader := func(t *testing.T, label string, values []float64) {
		t.Logf("\n=== %s ===", label)
		t.Logf("数据:   %v", floatSliceToString(values))
		indices := make([]int, len(values))
		for i := range indices {
			indices[i] = i
		}
		t.Logf("索引:   %v", intSliceToString(indices))
		t.Log("")
	}

	// 辅助函数：格式化输出结果
	printResult := func(t *testing.T, desc string, result PeaksResult, data []float64) {
		t.Logf("【%s】", desc)
		t.Logf("主趋势极值索引: %v", intSliceToString(result.Peaks))
		t.Logf("主趋势极值数值: %v", floatSliceToString(dataFromIndices(data, result.Peaks)))
		t.Logf("异常突破点索引: %v", intSliceToString(result.Breakouts))
		t.Log("")
	}

	// ========== 波峰测试 ==========
	printHeader(t, "📈 波峰检测数据", data)

	modes := func(left, right SearchMode) SideModes {
		return SideModes{Left: left, Right: right}
	}

	// 模式1：左侧找拐点，右侧保趋势
	result1 := FindExtremesWithBreakouts(data, 0, len(data), modes(FindInflection, PreserveTrend), ExtremePeak)
	printResult(t, "波峰 | 左侧: FindInflection | 右侧: PreserveTrend", result1, data)

	// 模式2：左侧保趋势，右侧找拐点
	result2 := FindExtremesWithBreakouts(data, 0, len(data), modes(PreserveTrend, FindInflection), ExtremePeak)
	printResult(t, "波峰 | 左侧: PreserveTrend | 右侧: FindInflection", result2, data)

	// 模式3：两侧都找拐点
	result3 := FindExtremesWithBreakouts(data, 0, len(data), modes(FindInflection, FindInflection), ExtremePeak)
	printResult(t, "波峰 | 左侧: FindInflection | 右侧: FindInflection", result3, data)

	// 模式4：两侧都保趋势
	result4 := FindExtremesWithBreakouts(data, 0, len(data), modes(PreserveTrend, PreserveTrend), ExtremePeak)
	printResult(t, "波峰 | 左侧: PreserveTrend | 右侧: PreserveTrend", result4, data)

	// ========== 波谷测试 ==========
	printHeader(t, "📉 波谷检测数据", lows)

	// 模式1：左侧找拐点，右侧保趋势
	valley1 := FindExtremesWithBreakouts(lows, 0, len(lows), modes(FindInflection, PreserveTrend), ExtremeTrough)
	printResult(t, "波谷 | 左侧: FindInflection | 右侧: PreserveTrend", valley1, lows)

	// 模式2：左侧保趋势，右侧找拐点
	valley2 := FindExtremesWithBreakouts(lows, 0, len(lows), modes(PreserveTrend, FindInflection), ExtremeTrough)
	printResult(t, "波谷 | 左侧: PreserveTrend | 右侧: FindInflection", valley2, lows)

	// 模式3：两侧都找拐点
	valley3 := FindExtremesWithBreakouts(lows, 0, len(lows), modes(FindInflection, FindInflection), ExtremeTrough)
	printResult(t, "波谷 | 左侧: FindInflection | 右侧: FindInflection", valley3, lows)

	// 模式4：两侧都保趋势
	valley4 := FindExtremesWithBreakouts(lows, 0, len(lows), modes(PreserveTrend, PreserveTrend), ExtremeTrough)
	printResult(t, "波谷 | 左侧: PreserveTrend | 右侧: PreserveTrend", valley4, lows)
}

// 辅助函数：将 float64 切片转为对齐字符串（固定宽度）
func floatSliceToString(f []float64) string {
	s := make([]string, len(f))
	for i, v := range f {
		s[i] = fmt.Sprintf("%.0f", v) // 整数格式
	}
	return padRight(strings.Join(s, " "), 4*len(s))
}

// 辅助函数：将 int 切片转为字符串并右填充
func intSliceToString(i []int) string {
	s := make([]string, len(i))
	for idx, v := range i {
		s[idx] = fmt.Sprintf("%d", v)
	}
	return padRight(strings.Join(s, " "), 4*len(s))
}

// 辅助函数：根据索引取值
func dataFromIndices(data []float64, indices []int) []float64 {
	var res []float64
	for _, i := range indices {
		if i >= 0 && i < len(data) {
			res = append(res, data[i])
		}
	}
	return res
}

// 辅助函数：右填充字符串（用于对齐）
func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}
