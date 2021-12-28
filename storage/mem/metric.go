package mem

import (
	"context"
	"fmt"

	"github.com/tyrylgin/collecter/model"
)

func (s *MemStore) GetAll(ctx context.Context) map[string]model.Metric {
	return s.metrics
}

func (s *MemStore) GetByName(ctx context.Context, name string) model.Metric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.metrics[name]; !ok {
		return nil
	}

	return s.metrics[name]
}

func (s *MemStore) Save(ctx context.Context, name string, metric model.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// if metric with same name registered, ignore
	if _, ok := s.metrics[name]; ok {
		return nil
	}

	switch metric.Type() {
	case model.MetricTypeCounter, model.MetricTypeGauge:
		s.metrics[name] = metric
	default:
		return fmt.Errorf("invalid type metric, support MetricTypeCounter/MetricTypeGauge")
	}

	return nil
}
