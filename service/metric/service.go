//go:generate mockgen -source=service.go -destination=./mock/metric.go -package=metricmock
package metric

import (
	"context"
	"errors"
	"fmt"

	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

type Service struct {
	metricStore storage.MetricStorer
}

func NewProcessor(storage storage.MetricStorer) *Service {
	srv := &Service{}
	srv.metricStore = storage
	return srv
}

type Processor interface {
	SetGauge(ctx context.Context, name string, value float64) error
	IncreaseCounter(ctx context.Context, name string, value int64) error
	GetAll(ctx context.Context) map[string]model.Metric
	Get(ctx context.Context, name string, metricType *model.MetricType) (model.Metric, error)
}

func (srv Service) SetGauge(ctx context.Context, name string, value float64) error {
	metric := srv.metricStore.GetByName(ctx, name)

	if metric == nil {
		newGauge := model.NewGauge()
		newGauge.Set(value)
		err := srv.metricStore.Save(ctx, name, newGauge.(model.Metric))
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

func (srv Service) IncreaseCounter(ctx context.Context, name string, value int64) error {
	metric := srv.metricStore.GetByName(ctx, name)

	if metric == nil {
		newGauge := model.NewCounter()
		newGauge.Increase(value)
		err := srv.metricStore.Save(ctx, name, newGauge.(model.Metric))
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

func (srv Service) GetAll(ctx context.Context) map[string]model.Metric {
	return srv.metricStore.GetAll(ctx)
}

func (srv Service) Get(ctx context.Context, name string, metricType *model.MetricType) (model.Metric, error) {
	metric := srv.metricStore.GetByName(ctx, name)

	if metric != nil && metricType != nil && metric.Type() != *metricType {
		return nil, fmt.Errorf("metric with name %s has different type %s", name, metric.Type())
	}

	return metric, nil
}
