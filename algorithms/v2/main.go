// plot.go
package main

import (
	"fmt"
	"image/color"
	"log"
	"os/exec"
	"runtime"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/quant1x/x/algorithms"
)

func main() {
	// 模拟 K 线数据
	high := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	low := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3}
	close := []float64{9, 11, 10, 12, 11, 13, 12, 14, 13, 18}

	// 构造 PriceSeries
	ps := &algorithms.PriceSeries{
		High:  high,
		Low:   low,
		Close: close,
	}

	// ✅ 真实调用你的算法函数
	sr := algorithms.FindSupportResistance(ps, 0, len(high))

	// ✅ 真实调用交易机会发现
	opportunities := algorithms.FindBreakoutOpportunities(ps, 0, len(high))

	// 创建绘图
	p := plot.New()
	p.Title.Text = "支撑/压力趋势分析（算法检测结果）"
	p.X.Label.Text = "时间"
	p.Y.Label.Text = "价格"

	// X 轴索引
	xs := make([]float64, len(high))
	for i := range xs {
		xs[i] = float64(i)
	}

	// ✅ 只绘制 High 和 Low 作为背景趋势线（不标点）
	highPoints := make(plotter.XYs, len(high))
	for i, h := range high {
		highPoints[i] = struct{ X, Y float64 }{X: xs[i], Y: h}
	}
	highLine, err := plotter.NewLine(highPoints)
	if err != nil {
		log.Fatal(err)
	}
	highLine.Color = color.RGBA{R: 255, G: 200, B: 180, A: 100} // 淡橙
	highLine.Width = vg.Points(0.5)
	p.Add(highLine)

	lowPoints := make(plotter.XYs, len(low))
	for i, l := range low {
		lowPoints[i] = struct{ X, Y float64 }{X: xs[i], Y: l}
	}
	lowLine, err := plotter.NewLine(lowPoints)
	if err != nil {
		log.Fatal(err)
	}
	lowLine.Color = color.RGBA{R: 180, G: 200, B: 255, A: 100} // 淡蓝
	lowLine.Width = vg.Points(0.5)
	p.Add(lowLine)

	// ✅ 标注算法检测出的主波峰（来自 FindSupportResistance）
	var resistancePoints plotter.XYs
	for _, idx := range sr.Resistance.Peaks {
		if idx >= 0 && idx < len(high) {
			resistancePoints = append(resistancePoints, struct{ X, Y float64 }{X: float64(idx), Y: high[idx]})
		}
	}
	if len(resistancePoints) > 0 {
		scatter, err := plotter.NewScatter(resistancePoints)
		if err != nil {
			log.Fatal(err)
		}
		scatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[0]         // Circle
		scatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红色
		scatter.GlyphStyle.Radius = vg.Points(6)
		p.Add(scatter)
		p.Legend.Add("主压力点 (High Peaks)", scatter)
	}

	// ✅ 标注算法检测出的主波谷（来自 FindSupportResistance）
	var supportPoints plotter.XYs
	for _, idx := range sr.Support.Peaks {
		if idx >= 0 && idx < len(low) {
			supportPoints = append(supportPoints, struct{ X, Y float64 }{X: float64(idx), Y: low[idx]})
		}
	}
	if len(supportPoints) > 0 {
		scatter, err := plotter.NewScatter(supportPoints)
		if err != nil {
			log.Fatal(err)
		}
		scatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[1]         // Square
		scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 200, B: 0, A: 255} // 绿色
		scatter.GlyphStyle.Radius = vg.Points(6)
		p.Add(scatter)
		p.Legend.Add("主支撑点 (Low Valleys)", scatter)
	}

	// ✅ 标注交易机会（来自递归分析）
	var oppPoints plotter.XYs
	for _, opp := range opportunities {
		var x float64
		var y float64
		if opp.Type == algorithms.ExtremeTrough {
			x = float64(opp.StartIdx)
			y = opp.Value
		} else if opp.Type == algorithms.ExtremePeak {
			x = float64(opp.StartIdx)
			y = opp.Value
		}
		oppPoints = append(oppPoints, struct{ X, Y float64 }{X: x, Y: y})
	}
	if len(oppPoints) > 0 {
		scatter, err := plotter.NewScatter(oppPoints)
		if err != nil {
			log.Fatal(err)
		}
		scatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[2]           // UpTriangle
		scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 255, B: 255, A: 255} // 青色
		scatter.GlyphStyle.Radius = vg.Points(8)
		p.Add(scatter)
		p.Legend.Add("交易机会", scatter)
	}

	// 保存图像
	if err := p.Save(10*vg.Inch, 7*vg.Inch, "trend_analysis.png"); err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println("✅ 图表已生成：trend_analysis.png")
	fmt.Printf("压力线被突破: %v\n", sr.Breakout.ResistanceBreak)
	fmt.Printf("支撑线被跌破: %v\n", sr.Breakout.SupportBreak)
	for _, opp := range opportunities {
		typ := "波峰"
		if opp.Type == algorithms.ExtremeTrough {
			typ = "波谷"
		}
		fmt.Printf("【交易机会】%s 在索引 %d, 值 %.2f\n", typ, opp.StartIdx, opp.Value)
	}

	// 自动打开（仅 Windows）
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "start", "trend_analysis.png")
		cmd.Start()
	}
}
