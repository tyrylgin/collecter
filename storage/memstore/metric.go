package memstore

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

func (s *MemStore) GetAll() map[string]model.Metric {
	return s.metrics
}

func (s *MemStore) Get(name string) model.Metric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.metrics[name]; !ok {
		return nil
	}

	return s.metrics[name]
}

func (s *MemStore) Save(name string, metric model.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// if metric with same name registered, ignore
	if _, ok := s.metrics[name]; ok {
		return nil
	}

	s.metrics[name] = metric

	return nil
}
