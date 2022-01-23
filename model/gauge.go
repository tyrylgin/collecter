package model

type Gauge struct {
	Value float64
}

func (g Gauge) Type() MetricType {
	return MetricTypeGauge
}
