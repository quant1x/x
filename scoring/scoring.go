package scoring

import (
	"errors"
	"fmt"
	"sync"
)

// DataProcessor 数据处理器接口
type DataProcessor interface {
	Process(interface{}) (float64, error)
}

// ScoreComponent 评分组件接口
type ScoreComponent interface {
	Name() string
	Weight() float64
	Score() (float64, error)
	SetSystem(*ScoringSystem)
	IsAutoWeight() bool
}

// BaseComponent 基础评分组件
type BaseComponent struct {
	name         string
	customWeight float64
	autoWeight   bool
	processor    DataProcessor
	rawData      interface{}
	system       *ScoringSystem
}

func (b *BaseComponent) Name() string {
	return b.name
}

func (b *BaseComponent) Weight() float64 {
	if b.IsAutoWeight() {
		return b.system.getAutoWeight()
	}
	return b.customWeight
}

func (b *BaseComponent) Score() (float64, error) {
	if b.processor == nil {
		return 0, errors.New("processor not defined")
	}
	return b.processor.Process(b.rawData)
}

func (b *BaseComponent) SetSystem(s *ScoringSystem) {
	b.system = s
}

func (b *BaseComponent) IsAutoWeight() bool {
	return b.autoWeight
}

// ScoringSystem 评分系统
type ScoringSystem struct {
	Components []ScoreComponent
	lock       sync.RWMutex
}

var instance *ScoringSystem
var once sync.Once

func GetScoringSystem() *ScoringSystem {
	once.Do(func() {
		instance = &ScoringSystem{
			Components: make([]ScoreComponent, 0),
		}
	})
	return instance
}

func (s *ScoringSystem) getAutoWeight() float64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	totalCustom := 0.0
	autoCount := 0

	for _, c := range s.Components {
		if !c.IsAutoWeight() {
			totalCustom += c.Weight()
		} else {
			autoCount++
		}
	}

	if totalCustom > 1 {
		return 0
	}

	remaining := 1 - totalCustom
	if autoCount == 0 {
		return 0
	}
	return remaining / float64(autoCount)
}

func (s *ScoringSystem) Register(c ScoreComponent) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, existing := range s.Components {
		if existing.Name() == c.Name() {
			return fmt.Errorf("component %s already exists", c.Name())
		}
	}

	c.SetSystem(s)
	s.Components = append(s.Components, c)
	return nil
}

func (s *ScoringSystem) CalculateTotalScore() (float64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(s.Components) == 0 {
		return 0, errors.New("no components registered")
	}

	total := 0.0
	for _, c := range s.Components {
		weight := c.Weight()
		score, err := c.Score()
		if err != nil {
			return 0, err
		}
		total += weight * score
	}

	if total > 100 {
		return 100, nil
	}
	return total, nil
}

// ComponentOption 组件选项
type ComponentOption func(*BaseComponent)

func NewComponent(name string, opts ...ComponentOption) ScoreComponent {
	c := &BaseComponent{
		name:       name,
		autoWeight: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := GetScoringSystem().Register(c); err != nil {
		panic(err)
	}
	return c
}

func WithCustomWeight(w float64) ComponentOption {
	return func(c *BaseComponent) {
		c.customWeight = w
		c.autoWeight = false
	}
}

func WithProcessor(p DataProcessor, data interface{}) ComponentOption {
	return func(c *BaseComponent) {
		c.processor = p
		c.rawData = data
	}
}
