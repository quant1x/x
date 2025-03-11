package scoring

import "fmt"

// RangeProcessor 示例处理器
type RangeProcessor struct {
	Min, Max float64
}

func (p *RangeProcessor) Process(data interface{}) (float64, error) {
	val, ok := data.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid data type")
	}

	if val < p.Min {
		return 0, nil
	}
	if val > p.Max {
		return 100, nil
	}
	return (val - p.Min) / (p.Max - p.Min) * 100, nil
}

// ThresholdProcessor 阀值处理器
type ThresholdProcessor struct {
	Threshold float64
}

func (p *ThresholdProcessor) Process(data interface{}) (float64, error) {
	val, ok := data.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid data type")
	}

	if val >= p.Threshold {
		return 100, nil
	}
	return val / p.Threshold * 100, nil
}
