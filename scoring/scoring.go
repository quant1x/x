package scoring

import (
	"errors"
)

// ScoreComponent 评分项结构体
type ScoreComponent struct {
	Name         string  // 维度名称
	CustomWeight float64 // 自定义权重（仅当IsAutoWeight=false时有效）
	Score        float64 // 得分（0-100）
	IsAutoWeight bool    // 是否自动分配权重
}

// ScoreCalculator 评分系统
type ScoreCalculator struct {
	components      []ScoreComponent
	computedWeights map[string]float64 // 存储计算后的实际权重
}

// NewScoreCalculator 创建评分系统实例
func NewScoreCalculator() *ScoreCalculator {
	return &ScoreCalculator{
		components:      make([]ScoreComponent, 0),
		computedWeights: make(map[string]float64),
	}
}

// AddComponent 添加评分组件（可选自定义权重）
func (s *ScoreCalculator) AddComponent(name string, score float64, customWeight ...float64) error {
	if score < 0 || score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	c := ScoreComponent{
		Name:         name,
		Score:        score,
		IsAutoWeight: true,
	}

	// 处理可选自定义权重
	if len(customWeight) > 0 {
		if len(customWeight) > 1 {
			return errors.New("too many weight parameters")
		}
		if customWeight[0] < 0 {
			return errors.New("weight cannot be negative")
		}
		c.CustomWeight = customWeight[0]
		c.IsAutoWeight = false
	}

	s.components = append(s.components, c)
	return nil
}

// CalculateWeightedScore 计算加权总分
func (s *ScoreCalculator) CalculateWeightedScore() (float64, error) {
	s.computedWeights = make(map[string]float64)

	// 计算自定义权重总和和自动分配项数量
	totalCustomWeight := 0.0
	autoCount := 0
	for _, c := range s.components {
		if !c.IsAutoWeight {
			totalCustomWeight += c.CustomWeight
		} else {
			autoCount++
		}
	}

	// 权重校验
	if totalCustomWeight > 1 {
		return 0, errors.New("total custom weight exceeds 1")
	}
	remainingWeight := 1 - totalCustomWeight
	if remainingWeight < 0 {
		return 0, errors.New("invalid weight distribution")
	}

	// 计算自动分配项的权重
	autoWeight := 0.0
	if autoCount > 0 {
		autoWeight = remainingWeight / float64(autoCount)
	}

	// 计算加权总分并记录实际权重
	totalScore := 0.0
	for _, c := range s.components {
		actualWeight := 0.0
		if !c.IsAutoWeight {
			actualWeight = c.CustomWeight
		} else {
			actualWeight = autoWeight
		}

		s.computedWeights[c.Name] = actualWeight
		totalScore += actualWeight * c.Score
	}

	return totalScore, nil
}

// WeightDistribution 获取权重分布
func (s *ScoreCalculator) WeightDistribution() map[string]float64 {
	return s.computedWeights
}
