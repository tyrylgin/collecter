package storage

import "github.com/tyrylgin/collecter/internal/metrics"

func StoreCounter(name string, value int64) {
	metric := metrics.Get(name)

	if metric == nil {
		newGauge := metrics.NewCounter()
		newGauge.Increase(value)
		metrics.Register(name, newGauge)
		return
	}

	counter := metric.(metrics.Counter)
	counter.Increase(value)
}
