package metrics

type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
)

type Metric interface {
	Type() MetricType
}
