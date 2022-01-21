package model

type Counter interface {
	GetDelta() int64
	IncreaseDelta(int64)
}

func NewCounter() Counter {
	return &DefaultCounter{0}
}

type DefaultCounter struct {
	delta int64
}

func (c *DefaultCounter) GetDelta() int64 {
	return c.delta
}

func (c *DefaultCounter) IncreaseDelta(i int64) {
	c.delta += i
}

func (c *DefaultCounter) Type() MetricType {
	return MetricTypeCounter
}
