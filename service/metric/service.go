//go:generate mockgen -source=service.go -destination=./mock/metric.go -package=metricmock
package metric

import (
	"errors"
	"fmt"

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
	Get(name string, metricType *model.MetricType) (model.Metric, error)
	GetAll() map[string]model.Metric
	IncreaseCounter(name string, value int64) error
	SetCounter(name string, value int64) error
	SetGauge(name string, value float64) error
}

func (s *Service) Get(name string, metricType *model.MetricType) (model.Metric, error) {
	metric := s.store.Get(name)

	if metric != nil && metricType != nil && metric.Type() != *metricType {
		return nil, fmt.Errorf("metric with name %s has different type %s", name, metric.Type())
	}

	return metric, nil
}

func (s *Service) GetAll() map[string]model.Metric {
	return s.store.GetAll()
}

func (s *Service) IncreaseCounter(name string, value int64) error {
	metric := s.store.Get(name)

	if metric == nil {
		newCounter := model.NewCounter()
		newCounter.Increase(value)
		err := s.store.Save(name, newCounter.(model.Metric))
		if err != nil {
			return fmt.Errorf("can't save new counter metric, %w", err)
		}
		return nil
	}

	if metric.Type() != model.MetricTypeCounter {
		return errors.New("the metric with same name already exist and has different type than MetricTypeCounter")
	}

	metric.(model.Counter).Increase(value)

	return nil
}

func (s *Service) SetCounter(name string, value int64) error {
	metric := s.store.Get(name)

	if metric == nil {
		newCounter := model.NewCounter()
		newCounter.Set(value)
		err := s.store.Save(name, newCounter.(model.Metric))
		if err != nil {
			return fmt.Errorf("can't save new counter metric, %w", err)
		}
		return nil
	}

	if metric.Type() != model.MetricTypeCounter {
		return errors.New("the metric with same name already exist and has different type than MetricTypeCounter")
	}

	counter := metric.(model.Counter)
	counter.Set(value)

	return nil
}

func (s *Service) SetGauge(name string, value float64) error {
	metric := s.store.Get(name)

	if metric == nil {
		newGauge := model.NewGauge()
		newGauge.Set(value)
		err := s.store.Save(name, newGauge.(model.Metric))
		if err != nil {
			return fmt.Errorf("can't save new gauge metric, %w", err)
		}
		return nil
	}

	if metric.Type() != model.MetricTypeGauge {
		return errors.New("the metric with same name already exist and has different type than MetricTypeGauge")
	}

	metric.(model.Gauge).Set(value)

	return nil
}
