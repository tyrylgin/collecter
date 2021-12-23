package storage

import "github.com/tyrylgin/collecter/internal/metrics"

func StoreGauge(name string, value float64) {
	metric := metrics.Get(name)

	if metric == nil {
		newGauge := metrics.NewGauge()
		newGauge.Set(value)
		metrics.Register(name, newGauge)
		return
	}

	gauge := metric.(metrics.Gauge)
	gauge.Set(value)
}
