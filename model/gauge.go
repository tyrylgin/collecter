package model

type Gauge interface {
	SetValue(float64)
	GetValue() float64
}

func NewGauge() Gauge {
	return &DefaultGauge{0}
}

type DefaultGauge struct {
	value float64
}

func (g *DefaultGauge) SetValue(v float64) {
	g.value = v
}

func (g *DefaultGauge) GetValue() float64 {
	return g.value
}

func (g *DefaultGauge) Type() MetricType {
	return MetricTypeGauge
}
