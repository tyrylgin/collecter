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

func (mm MetricMap) MarshalJSON() ([]byte, error) {
	type metricAlias struct {
		Metric Metric
		Type   string
	}
	metricAliasMap := map[string]metricAlias{}

	for name, metric := range mm {
		metricAliasMap[name] = metricAlias{
			Metric: metric,
			Type:   string(metric.Type()),
		}
	}

	return json.Marshal(metricAliasMap)
}

func (mm *MetricMap) UnmarshalJSON(data []byte) error {
	type metricAlias struct {
		Metric json.RawMessage
		Type   string
	}
	metricMapAlias := map[string]metricAlias{}
	if err := json.Unmarshal(data, &metricMapAlias); err != nil {
		return err
	}

	metrics := MetricMap{}
	for name, rawMetric := range metricMapAlias {
		switch rawMetric.Type {
		case string(MetricTypeCounter):
			counter := &Counter{}
			if err := json.Unmarshal(rawMetric.Metric, counter); err != nil {
				return fmt.Errorf("value of type (%T): unmarshal: %w", counter, err)
			}
			metrics[name] = counter
		case string(MetricTypeGauge):
			gauge := &Gauge{}
			if err := json.Unmarshal(rawMetric.Metric, gauge); err != nil {
				return fmt.Errorf("value of type (%T): unmarshal: %w", gauge, err)
			}
			metrics[name] = gauge
		}
	}

	*mm = metrics

	return nil
}
