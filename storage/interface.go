//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"

	"github.com/tyrylgin/collecter/model"
)

type MetricStorer interface {
	GetAll(ctx context.Context) map[string]model.Metric
	GetByName(ctx context.Context, name string) model.Metric
	Save(ctx context.Context, name string, metric model.Metric) error
}
