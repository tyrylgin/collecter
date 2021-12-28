package model

import "fmt"

type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

type Metric interface {
	Type() MetricType
}

func (t MetricType) Validate() error {
	switch t {
	case MetricTypeCounter, MetricTypeGauge:
		return nil
	default:
		return fmt.Errorf("unknown MetricType: %s", t)
	}
}
