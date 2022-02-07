//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"github.com/tyrylgin/collecter/model"
)

type MetricStorer interface {
	Get(name string) model.Metric
	GetAll() model.MetricMap
	Save(name string, metric model.Metric) error
	SaveAll(metrics model.MetricMap) error
}
