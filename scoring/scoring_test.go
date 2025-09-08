package scoring

import (
	"fmt"
	"testing"
)

func TestScoreCalculator(t *testing.T) {
	// 创建评分系统实例
	system := NewScoreCalculator()

	// 添加评分项（混合自定义和自动权重）
	system.AddComponent("用户体验", 80)     // 自动权重
	system.AddComponent("性能", 90, 0.3)  // 自定义权重0.3
	system.AddComponent("安全性", 70, 0.5) // 自定义权重0.5
	system.AddComponent("兼容性", 85)      // 自动权重

	// 计算最终得分
	score, err := system.CalculateWeightedScore()
	if err != nil {
		fmt.Println("计算失败:", err)
		return
	}

	// 输出结果
	fmt.Printf("综合得分: %.2f\n", score)
	fmt.Println("\n权重分布：")
	for name, weight := range system.WeightDistribution() {
		fmt.Printf("%s: %.1f%%\n", name, weight*100)
	}
}
