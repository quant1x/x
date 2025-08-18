package algorithms

//func TestWavesBasic1(t *testing.T) {
//	highs := []float64{10, 1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
//
//	t.Log("数据:", highs)
//	t.Log("索引: [0 1 2 3 4 5 6 7 8 9 10 11 12 13]")
//	t.Log("")
//
//	result1 := FindPeaksWithBreakouts(highs, 0, 14, FindInflection)
//	t.Log("【FindInflection 模式】")
//	t.Log("主趋势波峰索引:", result1.Peaks)
//	t.Log("主趋势波峰值:  ", dataFromIndices(highs, result1.Peaks))
//	t.Log("异常突破点索引:", result1.Breakouts)
//	t.Log("")
//
//	result2 := FindPeaksWithBreakouts(highs, 0, 14, PreserveTrend)
//	t.Log("【PreserveTrend 模式】")
//	t.Log("主趋势波峰索引:", result2.Peaks)
//	t.Log("主趋势波峰值:  ", dataFromIndices(highs, result2.Peaks))
//	t.Log("异常突破点索引:", result2.Breakouts)
//
//	//// 断言
//	//expected1 := []int{0, 2, 4, 6, 12}
//	//if !equal(result1.Peaks, expected1) {
//	//	t.Errorf("FindInflection: 期望 %v, 实际 %v", expected1, result1.Peaks)
//	//}
//	//
//	//expected2 := []int{0, 2, 12}
//	//if !equal(result2.Peaks, expected2) {
//	//	t.Errorf("PreserveTrend: 期望 %v, 实际 %v", expected2, result2.Peaks)
//	//}
//	fmt.Println("----------")
//
//	lows := []float64{5, 3, 6, 2, 4, 1, 7, 1, 8}
//	result := FindValleysWithBreakouts(lows, 0, len(lows), FindInflection)
//	// result.Peaks: 主要波谷位置（如 3, 5, 7 等）
//	// result.Breakouts: 异常抬高的波谷点
//	fmt.Println(result.Peaks)
//}
