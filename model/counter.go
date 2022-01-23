package model

type Counter struct {
	Delta int64
}

func (c Counter) Type() MetricType {
	return MetricTypeCounter
}
