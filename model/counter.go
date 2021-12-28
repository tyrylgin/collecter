package model

type Counter interface {
	Count() int64
	Increase(int64)
}

func NewCounter() Counter {
	return &DefaultCounter{0}
}

type DefaultCounter struct {
	count int64
}

func (c *DefaultCounter) Count() int64 {
	return c.count
}

func (c *DefaultCounter) Increase(i int64) {
	c.count += i
}

func (c *DefaultCounter) Type() MetricType {
	return MetricTypeCounter
}
