package main

import (
	"fmt"
	"github.com/quant1x/x/scoring"
)

func main() {
	// 创建评分组件
	scoring.NewComponent("性能",
		scoring.WithProcessor(&scoring.RangeProcessor{Min: 0, Max: 200}, 150.0),
		scoring.WithCustomWeight(0.3),
	)

	scoring.NewComponent("安全",
		scoring.WithProcessor(&scoring.ThresholdProcessor{Threshold: 90}, 85.0),
	)

	scoring.NewComponent("可用性",
		scoring.WithProcessor(&scoring.RangeProcessor{Min: 0.9, Max: 1.0}, 0.95),
		scoring.WithCustomWeight(0.2),
	)

	// 计算总分
	system := scoring.GetScoringSystem()
	score, err := system.CalculateTotalScore()
	if err != nil {
		panic(err)
	}

	fmt.Printf("系统综合得分: %.2f\n", score)

	// 显示权重分布
	fmt.Println("\n权重分布：")
	for _, c := range system.Components {
		fmt.Printf("%-10s : %.1f%%\n", c.Name(), c.Weight()*100)
	}
}
