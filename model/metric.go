package model

import (
	"encoding/json"
	"fmt"
)

type (
	MetricType string

	Metric interface {
		Type() MetricType
	}

	MetricMap map[string]Metric
)

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

func (t MetricType) Validate() error {
	switch t {
	case MetricTypeCounter, MetricTypeGauge:
		return nil
	default:
		return fmt.Errorf("unknown MetricType: %s", t)
	}
}
