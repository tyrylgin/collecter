package memstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

var _ storage.MetricStorer = (*MemStore)(nil)

type MemStore struct {
	metrics      model.MetricMap
	mutex        sync.RWMutex
	file         *os.File
	isSyncBackup bool
}

func NewStorage() MemStore {
	return MemStore{metrics: model.MetricMap{}}
}

func (s *MemStore) GetAll() model.MetricMap {
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

	s.metrics[name] = metric

	return nil
}
