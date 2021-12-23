package metrics

import (
	"log"
	"sync"
)

type Storage struct {
	metrics map[string]Metric
	mutex   sync.RWMutex
}

var DefaultStorage = &Storage{metrics: make(map[string]Metric)}

func Register(name string, i interface{}) {
	DefaultStorage.mutex.Lock()
	defer DefaultStorage.mutex.Unlock()

	// If metric with same name registered, ignore
	if _, ok := DefaultStorage.metrics[name]; ok {
		return
	}

	switch i.(type) {
	case Gauge, Counter:
		DefaultStorage.metrics[name] = i.(Metric)
	default:
		log.Fatal("unknown metric type, allowed MetricTypeCounter and MetricTypeGauge")
		return
	}
}

func Get(name string) Metric {
	DefaultStorage.mutex.RLock()
	defer DefaultStorage.mutex.RUnlock()
	return DefaultStorage.metrics[name]
}

func Each(f func(string, Metric)) {
	DefaultStorage.mutex.RLock()
	defer DefaultStorage.mutex.RUnlock()

	for name, metric := range DefaultStorage.metrics {
		f(name, metric)
	}
}
