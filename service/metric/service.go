//go:generate mockgen -source=service.go -destination=./mock/metric.go -package=metricmock
package metric

import (
	"github.com/pkg/errors"
	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

type Service struct {
	store storage.MetricStorer
}

func NewProcessor(store storage.MetricStorer) *Service {
	return &Service{
		store: store,
	}
}

type Processor interface {
	Get(name string, metricType model.MetricType) (model.Metric, error)
	GetAll() map[string]model.Metric
	IncreaseCounter(name string, value int64) error
	SetGauge(name string, value float64) error
	SetMetrics(metrics map[string]model.Metric) error
}

func (s *Service) Get(name string, metricType model.MetricType) (model.Metric, error) {
	metric := s.store.Get(name)

	if metric != nil && metric.Type() != metricType {
		return nil, errors.Errorf("metric with name %s has different type %s", name, metric.Type())
	}

	return metric, nil
}

func (s *Service) GetAll() map[string]model.Metric {
	return s.store.GetAll()
}

func (s *Service) IncreaseCounter(name string, value int64) error {
	metric, ok := s.store.Get(name).(model.Counter)

	if !ok {
		metric = model.Counter{}
	}

	metric.Delta += value
	if err := s.store.Save(name, metric); err != nil {
		return errors.Errorf("can't save counter metric, %v", err)
	}

	return nil
}

func (s *Service) SetGauge(name string, value float64) error {
	metric, ok := s.store.Get(name).(model.Gauge)

	if !ok {
		metric = model.Gauge{}
	}

	metric.Value = value
	if err := s.store.Save(name, metric); err != nil {
		return errors.Errorf("can't save counter metric, %v", err)
	}

	return nil
}

func (s *Service) SetMetrics(metrics map[string]model.Metric) error {
	if err := s.store.SaveAll(metrics); err != nil {
		return errors.Errorf("can't save metrics, %v", err)
	}

	return nil
}
