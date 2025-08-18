package algorithms

import (
	"fmt"
	"math"
	"slices"
	"sort"
)

const (
	WaveFloatEps = 1e-9
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

func (r PeaksResult) HasBreakouts() bool {
	return len(r.Breakouts) > 0
}

func (r PeaksResult) Count() int {
	return len(r.Peaks)
}

func (r PeaksResult) FirstPeak() int {
	if len(r.Peaks) == 0 {
		return -1
	}
	return r.Peaks[0]
}

func (r PeaksResult) LastPeak() int {
	if len(r.Peaks) == 0 {
		return -1
	}
	return r.Peaks[len(r.Peaks)-1]
}

func (r PeaksResult) Values(data []float64) []float64 {
	var vals []float64
	for _, i := range r.Peaks {
		if i >= 0 && i < len(data) {
			vals = append(vals, data[i])
		}
	}
	return vals
}

func (r PeaksResult) IsEmpty() bool {
	return len(r.Peaks) == 0
}

// SideModes 允许为左侧和右侧自由段独立设置检测模式
type SideModes struct {
	Left  SearchMode // 第一个主峰/主谷左侧使用的模式
	Right SearchMode // 最后一个主峰/主谷右侧使用的模式
}

// checkAndAppend 判断当前点是否符合趋势（非递增/非递减）
func checkAndAppend(
	data []float64,
	currIdx int,
	valid *[]int,
	breakouts *[]int,
	shouldIncrease bool,
) {
	if len(*valid) == 0 {
		*valid = append(*valid, currIdx)
		return
	}

	lastVal := data[(*valid)[len(*valid)-1]]
	currVal := data[currIdx]

	if shouldIncrease {
		if currVal >= lastVal {
			*valid = append(*valid, currIdx)
		} else {
			*breakouts = append(*breakouts, currIdx)
		}
	} else {
		if currVal <= lastVal {
			*valid = append(*valid, currIdx)
		} else {
			*breakouts = append(*breakouts, currIdx)
		}
	}
}

// Option 配置选项
type Option func(*config)

type config struct {
	eps float64
}

// WithEpsilon 设置浮点比较精度
func WithEpsilon(eps float64) Option {
	return func(c *config) {
		c.eps = eps
	}
}

// FindExtremesWithBreakouts 使用外部提供的 extremes 列表，或降级为自动提取，
// 在 [start, end) 区间分析波峰或波谷。
//
// 核心思想：以“全局最值的所有位置”为锚点，将序列分为：
//
//	左自由段 | 主峰区 | 右自由段
//	对左右自由段中的候选极值点进行趋势合规性审查
//
// 参数：
//   - data: 时间序列
//   - extremes: 外部提供的候选转折点（可为 nil 或空，此时自动提取局部极值）
//   - start, end: 分析区间 [start, end)
//   - modes: 左右自由段的审查模式
//   - direction: ExtremePeak（波峰）或 ExtremeTrough（波谷）
//   - opts: 可选配置（如 WithEpsilon）
func FindExtremesWithBreakouts(
	data []float64,
	extremes []int,
	start, end int,
	modes SideModes,
	direction ExtremeType,
	opts ...Option,
) PeaksResult {
	result := PeaksResult{}

	// 应用选项
	cfg := config{eps: WaveFloatEps}
	for _, opt := range opts {
		opt(&cfg)
	}
	eps := cfg.eps

	// 防护
	if data == nil || len(data) == 0 ||
		start < 0 || end > len(data) || start >= end || len(data) < 3 {
		return result
	}

	// 🔽 降级处理：如果 extremes 为空，则自动提取局部极值
	if extremes == nil || len(extremes) == 0 {
		tempExtremes := findLocalExtremesIn(data, start, end, direction)
		if len(tempExtremes) == 0 {
			// 即使降级也无法提取极值 → 只保留主极值点
			extremes = nil // 清空，后续跳过审查
		} else {
			extremes = tempExtremes
		}
	}

	// =============================================
	// 🔑 第一步：找全局最值 mainVal（在整个 [start, end) 区间）
	// =============================================
	var mainVal float64
	if direction == ExtremePeak {
		mainVal = data[start]
		for i := start; i < end; i++ {
			if data[i] > mainVal {
				mainVal = data[i]
			}
		}
	} else {
		mainVal = data[start]
		for i := start; i < end; i++ {
			if data[i] < mainVal {
				mainVal = data[i]
			}
		}
	}

	// =============================================
	// 🔑 第二步：收集所有等于 mainVal 的点 → 主极值点（值驱动，非结构驱动）
	// =============================================
	var majorExtremes []int
	for i := start; i < end; i++ {
		if math.Abs(data[i]-mainVal) < eps {
			majorExtremes = append(majorExtremes, i)
		}
	}

	// 安全兜底
	if len(majorExtremes) == 0 {
		return result
	}

	sort.Ints(majorExtremes)
	firstMajor := majorExtremes[0]
	lastMajor := majorExtremes[len(majorExtremes)-1]

	// =============================================
	// 🌐 第三步：分片治理
	//   - 左自由段: [start, firstMajor) → 审查趋势
	//   - 主峰区: [firstMajor, lastMajor] → 全部保留
	//   - 右自由段: (lastMajor, end) → 审查趋势
	// =============================================
	var results []int
	var breakouts []int

	// 只有在有 extremes 的情况下才审查
	if extremes != nil && len(extremes) > 0 {
		// 1. 处理左侧自由段
		processExtremeSegmentWithEps(
			data, start, firstMajor,
			extremes, modes.Left,
			&results, &breakouts,
			mainVal, direction, SideLeft,
			eps,
		)

		// 2. 处理右侧自由段
		processExtremeSegmentWithEps(
			data, lastMajor+1, end,
			extremes, modes.Right,
			&results, &breakouts,
			mainVal, direction, SideRight,
			eps,
		)
	}

	// 3. 加入主峰区所有点
	results = append(results, majorExtremes...)

	// 4. 排序输出
	sort.Ints(results)
	sort.Ints(breakouts)

	result.Peaks = results
	result.Breakouts = breakouts
	return result
}

// processExtremeSegmentWithEps 支持自定义 epsilon
func processExtremeSegmentWithEps(
	data []float64,
	segStart, segEnd int,
	extremes []int,
	mode SearchMode,
	results *[]int,
	breakouts *[]int,
	mainVal float64,
	direction ExtremeType,
	side SegmentSide,
	eps float64,
) {
	if segStart < 0 || segEnd > len(data) || segStart >= segEnd {
		return
	}

	var segExtremes []int
	for _, idx := range extremes {
		if idx >= segStart && idx < segEnd && math.Abs(data[idx]-mainVal) >= eps {
			segExtremes = append(segExtremes, idx)
		}
	}

	if len(segExtremes) == 0 {
		return
	}

	var valid []int
	var increasing bool // true: 非递减；false: 非递增
	var reverseOrder bool

	switch {
	case mode == FindInflection && side == SideLeft && direction == ExtremePeak:
		reverseOrder = true
		increasing = false
	case mode == FindInflection && side == SideRight && direction == ExtremePeak:
		reverseOrder = false
		increasing = false
	case mode == PreserveTrend && side == SideLeft && direction == ExtremePeak:
		reverseOrder = false
		increasing = true
	case mode == PreserveTrend && side == SideRight && direction == ExtremePeak:
		reverseOrder = true
		increasing = true

	case mode == FindInflection && side == SideLeft && direction == ExtremeTrough:
		reverseOrder = true
		increasing = true
	case mode == FindInflection && side == SideRight && direction == ExtremeTrough:
		reverseOrder = false
		increasing = true
	case mode == PreserveTrend && side == SideLeft && direction == ExtremeTrough:
		reverseOrder = false
		increasing = false
	case mode == PreserveTrend && side == SideRight && direction == ExtremeTrough:
		reverseOrder = true
		increasing = false
	}

	indices := segExtremes
	if reverseOrder {
		for i := len(segExtremes) - 1; i >= 0; i-- {
			checkAndAppend(data, segExtremes[i], &valid, breakouts, increasing)
		}
	} else {
		for _, idx := range indices {
			checkAndAppend(data, idx, &valid, breakouts, increasing)
		}
	}

	if reverseOrder {
		slices.Reverse(valid)
	}

	*results = append(*results, valid...)
}

func findLocalExtremesIn(data []float64, start, end int, direction ExtremeType) []int {
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

	return extremes
}

// 默认模式：左侧保趋势，右侧找拐点
var defaultPeakModes = SideModes{
	Left:  PreserveTrend,  // 左侧既成事实，保留
	Right: FindInflection, // 右侧审查是否破坏趋势
}

var defaultValleyModes = SideModes{
	Left:  PreserveTrend,
	Right: FindInflection,
}

// FindPeaks 在 [start, end) 区间检测波峰
//
// 默认行为：
//   - extremes: nil → 自动提取局部极值
//   - modes:    根据主峰位置智能选择左侧模式
//   - eps:      使用默认 floatEps
func FindPeaks(data []float64, start, end int, modes SideModes, opts ...Option) PeaksResult {
	// 1. 如果用户已指定模式，直接使用
	if modes.Left == 0 && modes.Right == 0 {
		modes = defaultPeakModes
	}

	// 2. 否则：智能推断模式
	return findPeaksWithAutoModes(data, start, end, opts...)
}

// findPeaksWithAutoModes 根据主峰位置自动选择模式
func findPeaksWithAutoModes(data []float64, start, end int, opts ...Option) PeaksResult {
	// 先提取局部极值（用于判断主峰位置）
	extremes := findLocalExtremesIn(data, start, end, ExtremePeak)
	if len(extremes) == 0 {
		// 无法提取极值 → 使用保守模式
		modes := SideModes{
			Left:  FindInflection,
			Right: PreserveTrend,
		}
		return FindExtremesWithBreakouts(data, nil, start, end, modes, ExtremePeak, opts...)
	}

	// 找主值
	mainVal := data[extremes[0]]
	for _, i := range extremes {
		if data[i] > mainVal {
			mainVal = data[i]
		}
	}

	// 应用选项
	cfg := config{eps: WaveFloatEps}
	for _, opt := range opts {
		opt(&cfg)
	}
	eps := cfg.eps

	// 收集主极值点
	var majorExtremes []int
	for _, i := range extremes {
		if math.Abs(data[i]-mainVal) < eps {
			majorExtremes = append(majorExtremes, i)
		}
	}
	if len(majorExtremes) == 0 {
		return PeaksResult{}
	}

	sort.Ints(majorExtremes)
	lastMajor := majorExtremes[len(majorExtremes)-1] // 最后一个主峰位置

	// 判断趋势：主峰是否在右半区？
	midPoint := start + int(float64(end-start)*0.6) // 60% 分位为界

	var leftMode SearchMode
	if lastMajor >= midPoint {
		// 主峰靠右 → 上涨趋势 → 左侧应保趋势
		leftMode = PreserveTrend
	} else {
		// 主峰靠左或居中 → 可能见顶 → 左侧找拐点
		leftMode = FindInflection
	}

	modes := SideModes{
		Left:  leftMode,
		Right: PreserveTrend, // 右侧默认保趋势
	}

	return FindExtremesWithBreakouts(data, nil, start, end, modes, ExtremePeak, opts...)
}

// FindValleys 智能模式：主谷靠右 → 下降趋势 → 左侧保趋势（下降）
func FindValleys(data []float64, start, end int, modes SideModes, opts ...Option) PeaksResult {
	if modes.Left == 0 && modes.Right == 0 {
		modes = defaultValleyModes
	}
	return findValleysWithAutoModes(data, start, end, opts...)
}

func findValleysWithAutoModes(data []float64, start, end int, opts ...Option) PeaksResult {
	extremes := findLocalExtremesIn(data, start, end, ExtremeTrough)
	if len(extremes) == 0 {
		modes := SideModes{
			Left:  FindInflection,
			Right: PreserveTrend,
		}
		return FindExtremesWithBreakouts(data, nil, start, end, modes, ExtremeTrough, opts...)
	}

	mainVal := data[extremes[0]]
	for _, i := range extremes {
		if data[i] < mainVal {
			mainVal = data[i]
		}
	}

	// 应用选项
	cfg := config{eps: WaveFloatEps}
	for _, opt := range opts {
		opt(&cfg)
	}
	eps := cfg.eps

	var majorExtremes []int
	for _, i := range extremes {
		if math.Abs(data[i]-mainVal) < eps {
			majorExtremes = append(majorExtremes, i)
		}
	}
	if len(majorExtremes) == 0 {
		return PeaksResult{}
	}

	sort.Ints(majorExtremes)
	lastMajor := majorExtremes[len(majorExtremes)-1]
	midPoint := start + int(float64(end-start)*0.6)

	var leftMode SearchMode
	if lastMajor >= midPoint {
		// 主谷靠右 → 持续下跌 → 左侧保趋势（下降）
		leftMode = PreserveTrend
	} else {
		// 主谷靠左 → 早期探底 → 左侧找拐点
		leftMode = FindInflection
	}

	modes := SideModes{
		Left:  leftMode,
		Right: PreserveTrend,
	}

	return FindExtremesWithBreakouts(data, nil, start, end, modes, ExtremeTrough, opts...)
}

// PriceSeries K线数据序列
type PriceSeries struct {
	High       []float64 // 最高价序列
	Low        []float64 // 最低价序列
	Close      []float64 // 收盘价（可选）
	Timestamps []int64   // 时间戳（可选）
}

// SupportResistance 保存支撑/压力趋势结果
type SupportResistance struct {
	Resistance PeaksResult // 压力线（来自 high 的波峰）
	Support    PeaksResult // 支撑线（来自 low 的波谷）
	Breakout   struct {
		ResistanceBreak bool // 压力线被突破
		SupportBreak    bool // 支撑线被跌破
		FirstBreakIdx   int  // 首次突破位置
	}
}

// FindSupportResistance 从 high/low 序列中提取支撑/压力趋势
func FindSupportResistance(ps *PriceSeries, start, end int) SupportResistance {
	var sr SupportResistance

	if end <= start {
		return sr
	}

	modes := SideModes{
		Left:  PreserveTrend,
		Right: FindInflection,
	}

	// 1. 找极值
	sr.Resistance = FindExtremesWithBreakouts(ps.High, nil, start, end, modes, ExtremePeak)
	sr.Support = FindExtremesWithBreakouts(ps.Low, nil, start, end, modes, ExtremeTrough)

	// 2. 获取“突破前”的主压力（即：最后一个主峰之前的最高主压力）
	var prevResistance float64 = math.Inf(-1)
	var lastPeakIdx int = -1
	if len(sr.Resistance.Peaks) > 0 {
		// 排序确保有序
		sort.Ints(sr.Resistance.Peaks)
		// 取最后一个主峰索引
		lastPeakIdx = sr.Resistance.Peaks[len(sr.Resistance.Peaks)-1]
		// 找出在 lastPeakIdx 之前的主峰中的最大值
		for _, idx := range sr.Resistance.Peaks {
			if idx < lastPeakIdx && ps.High[idx] > prevResistance {
				prevResistance = ps.High[idx]
			}
		}
		// 如果没有之前的主峰，用全局次高？
		if prevResistance == math.Inf(-1) {
			prevResistance = 0 // 或设为最小值
		}
	}

	// 3. 判断是否突破：当前 High 是否 > 之前主压力
	for i := start; i < end; i++ {
		// 如果该点是主峰，且其值 > 之前主压力，且不是第一个主峰 → 视为突破
		if contains(sr.Resistance.Peaks, i) && ps.High[i] > prevResistance {
			sr.Breakout.ResistanceBreak = true
			if sr.Breakout.FirstBreakIdx == 0 || i < sr.Breakout.FirstBreakIdx {
				sr.Breakout.FirstBreakIdx = i
			}
		}
	}

	return sr
}

// 辅助函数
func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// getMaxPeakValue 获取主压力值（最后一个主峰）
func getMaxPeakValue(high []float64, peaks []int) float64 {
	if len(peaks) == 0 {
		return math.Inf(-1)
	}
	lastPeak := peaks[len(peaks)-1]
	return high[lastPeak]
}

// getMinValleyValue 获取主支撑值（最后一个主谷）
func getMinValleyValue(low []float64, peaks []int) float64 {
	if len(peaks) == 0 {
		return math.Inf(1)
	}
	lastValley := peaks[len(peaks)-1]
	return low[lastValley]
}

// isInFuture 判断索引是否在最后一个主极值之后
func isInFuture(idx int, peaks []int) bool {
	if len(peaks) == 0 {
		return true
	}
	lastPeak := peaks[len(peaks)-1]
	return idx > lastPeak
}

// TradeOpportunity 表示一个可交易的局部趋势机会
// 通常由趋势破坏后，递归发现的新局部主趋势构成
type TradeOpportunity struct {
	Type      ExtremeType    // 机会类型：ExtremePeak（高点卖出/做空）或 ExtremeTrough（低点买入/做多）
	StartIdx  int            // 该机会分析的起始索引
	EndIdx    int            // 该机会分析的结束索引
	Peaks     []int          // 合规的趋势点（主趋势+自由段合规点）
	Breakouts []int          // 异常突破点（破坏趋势的点）
	Value     float64        // 主极值点的值（价格）
	Direction TrendDirection // 趋势方向（可选）
}

// TrendDirection 趋势方向（辅助判断）
type TrendDirection int

const (
	TrendUnknown  TrendDirection = iota
	TrendUpward                  // 上升趋势（如：波谷后回升）
	TrendDownward                // 下降趋势（如：波峰后回落）
)

// String 方法（便于日志输出）
func (t TrendDirection) String() string {
	switch t {
	case TrendUpward:
		return "Upward"
	case TrendDownward:
		return "Downward"
	default:
		return "Unknown"
	}
}

// String 方法
func (to TradeOpportunity) String() string {
	typ := "Peak"
	if to.Type == ExtremeTrough {
		typ = "Trough"
	}
	return fmt.Sprintf("TradeOpportunity{Type: %s, Start: %d, End: %d, Value: %.3f, Direction: %s, Peaks: %v, Breakouts: %v}",
		typ, to.StartIdx, to.EndIdx, to.Value, to.Direction, to.Peaks, to.Breakouts)
}

// FindBreakoutOpportunities 在突破后寻找新趋势
func FindBreakoutOpportunities(ps *PriceSeries, start, end int) []TradeOpportunity {
	sr := FindSupportResistance(ps, start, end)
	var opportunities []TradeOpportunity

	if sr.Breakout.FirstBreakIdx != 0 {
		// 从突破点开始，找新趋势
		subStart := sr.Breakout.FirstBreakIdx
		subEnd := end

		subModes := SideModes{
			Left:  PreserveTrend,
			Right: PreserveTrend,
		}

		// 突破压力 → 找新上升趋势（low 中找波谷结束回调）
		if sr.Breakout.ResistanceBreak {
			subValley := FindExtremesWithBreakouts(ps.Low, nil, subStart, subEnd, subModes, ExtremeTrough)
			if len(subValley.Peaks) > 0 {
				last := subValley.Peaks[len(subValley.Peaks)-1]
				opportunities = append(opportunities, TradeOpportunity{
					Type:      ExtremeTrough,
					StartIdx:  subStart,
					EndIdx:    subEnd,
					Peaks:     subValley.Peaks,
					Breakouts: subValley.Breakouts,
					Value:     ps.Low[last],
				})
			}
		}

		// 跌破支撑 → 找新下降趋势（high 中找波峰反弹结束）
		if sr.Breakout.SupportBreak {
			subPeak := FindExtremesWithBreakouts(ps.High, nil, subStart, subEnd, subModes, ExtremePeak)
			if len(subPeak.Peaks) > 0 {
				last := subPeak.Peaks[len(subPeak.Peaks)-1]
				opportunities = append(opportunities, TradeOpportunity{
					Type:      ExtremePeak,
					StartIdx:  subStart,
					EndIdx:    subEnd,
					Peaks:     subPeak.Peaks,
					Breakouts: subPeak.Breakouts,
					Value:     ps.High[last],
				})
			}
		}
	}

	return opportunities
}
