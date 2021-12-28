package mem

import (
	"sync"

	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

var _ storage.MetricStorer = (*MemStore)(nil)

type MemStore struct {
	metrics map[string]model.Metric
	mutex   sync.RWMutex
}

func NewStorage() MemStore {
	return MemStore{metrics: make(map[string]model.Metric)}
}
