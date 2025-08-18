// peaks_test.go
package algorithms

import (
	"fmt"
	"strings"
	"testing"
)

func TestWavesV2(t *testing.T) {
	highs := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	//highs := []float64{1, 5, 5, 4, 5, 5, 1}
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
	printResult := func(t *testing.T, desc string, result PeaksResult, highs []float64) {
		t.Logf("【%s】", desc)
		t.Logf("主趋势极值索引: %v", intSliceToString(result.Peaks))
		t.Logf("主趋势极值数值: %v", floatSliceToString(dataFromIndices(highs, result.Peaks)))
		t.Logf("异常突破点索引: %v", intSliceToString(result.Breakouts))
		t.Log("")
	}

	// ========== 波峰测试 ==========
	printHeader(t, "📈 波峰检测数据", highs)

	modes := func(left, right SearchMode) SideModes {
		return SideModes{Left: left, Right: right}
	}

	// 模式1：左侧找拐点，右侧保趋势
	result1 := FindExtremesWithBreakouts(highs, nil, 0, len(highs), modes(FindInflection, PreserveTrend), ExtremePeak)
	printResult(t, "波峰 | 左侧: FindInflection | 右侧: PreserveTrend", result1, highs)

	// 模式2：左侧保趋势，右侧找拐点
	result2 := FindExtremesWithBreakouts(highs, nil, 0, len(highs), modes(PreserveTrend, FindInflection), ExtremePeak)
	printResult(t, "波峰 | 左侧: PreserveTrend | 右侧: FindInflection", result2, highs)

	// 模式3：两侧都找拐点
	result3 := FindExtremesWithBreakouts(highs, nil, 0, len(highs), modes(FindInflection, FindInflection), ExtremePeak)
	printResult(t, "波峰 | 左侧: FindInflection | 右侧: FindInflection", result3, highs)

	// 模式4：两侧都保趋势
	result4 := FindExtremesWithBreakouts(highs, nil, 0, len(highs), modes(PreserveTrend, PreserveTrend), ExtremePeak)
	printResult(t, "波峰 | 左侧: PreserveTrend | 右侧: PreserveTrend", result4, highs)

	// ========== 波谷测试 ==========
	printHeader(t, "📉 波谷检测数据", lows)

	// 模式1：左侧找拐点，右侧保趋势
	valley1 := FindExtremesWithBreakouts(lows, nil, 0, len(lows), modes(FindInflection, PreserveTrend), ExtremeTrough)
	printResult(t, "波谷 | 左侧: FindInflection | 右侧: PreserveTrend", valley1, lows)

	// 模式2：左侧保趋势，右侧找拐点
	valley2 := FindExtremesWithBreakouts(lows, nil, 0, len(lows), modes(PreserveTrend, FindInflection), ExtremeTrough)
	printResult(t, "波谷 | 左侧: PreserveTrend | 右侧: FindInflection", valley2, lows)

	// 模式3：两侧都找拐点
	valley3 := FindExtremesWithBreakouts(lows, nil, 0, len(lows), modes(FindInflection, FindInflection), ExtremeTrough)
	printResult(t, "波谷 | 左侧: FindInflection | 右侧: FindInflection", valley3, lows)

	// 模式4：两侧都保趋势
	valley4 := FindExtremesWithBreakouts(lows, nil, 0, len(lows), modes(PreserveTrend, PreserveTrend), ExtremeTrough)
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

func TestBreakoutOpportunities(t *testing.T) {
	ps := &PriceSeries{
		High:  []float64{10, 12, 11, 13, 12, 14, 13, 15, 14, 20},
		Low:   []float64{8, 9, 8, 10, 9, 11, 10, 12, 11, 13},
		Close: []float64{9, 11, 10, 12, 11, 13, 12, 14, 13, 18},
	}

	sr := FindSupportResistance(ps, 0, len(ps.High))

	// 🔍 调试输出
	t.Log("Resistance.Peaks:", sr.Resistance.Peaks)
	t.Log("Resistance.Breakouts:", sr.Resistance.Breakouts)
	for i, h := range ps.High {
		t.Log(fmt.Sprintf("High[%d] = %.2f", i, h))
	}

	fmt.Printf("压力线被突破: %v\n", sr.Breakout.ResistanceBreak)
	fmt.Printf("支撑线被跌破: %v\n", sr.Breakout.SupportBreak)

	opportunities := FindBreakoutOpportunities(ps, 0, len(ps.High))
	for _, opp := range opportunities {
		fmt.Printf("新机会: %v, 位置: %d, 值: %.2f\n", opp.Type, opp.StartIdx, opp.Value)
	}
	for _, opp := range opportunities {
		if opp.Type == ExtremeTrough {
			fmt.Println("【买入信号】在", opp.Value, "找到支撑")
		} else if opp.Type == ExtremePeak {
			fmt.Println("【卖出信号】在", opp.Value, "遇到压力")
		}
	}
}
