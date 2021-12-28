package model

type Gauge interface {
	Set(float64)
	Value() float64
}

func NewGauge() Gauge {
	return &DefaultGauge{0}
}

type DefaultGauge struct {
	value float64
}

func (g *DefaultGauge) Set(v float64) {
	g.value = v
}

func (g *DefaultGauge) Value() float64 {
	return g.value
}

func (g *DefaultGauge) Type() MetricType {
	return MetricTypeGauge
}
