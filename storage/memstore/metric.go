package memstore

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
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

func NewStorage() *MemStore {
	return &MemStore{metrics: model.MetricMap{}}
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

	if s.isSyncBackup {
		if err := s.dropToFile(); err != nil {
			return errors.Wrap(err, "failed to backup metrics to file in sync mode")
		}
	}

	return nil
}

func (s *MemStore) SaveAll(metrics model.MetricMap) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.metrics = metrics

	if s.isSyncBackup {
		if err := s.dropToFile(); err != nil {
			return errors.Wrap(err, "failed to backup metrics to file in sync mode")
		}
	}

	return nil
}

func (s *MemStore) WithFileBackup(ctx context.Context, fileName string, storeInterval time.Duration, isRestore bool) (err error) {
	if fileName == "" {
		return nil
	}

	s.file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errors.Wrap(err, "can't open file for memstore")
	}

	if isRestore {
		if err := s.restoreFromFile(); err != nil {
			log.Printf("failed to load metrics backup from file; %v", err)
		}
	}

	if storeInterval > 0 {
		go func() {
			ticker := time.NewTicker(storeInterval)
			for {
				<-ticker.C
				if err := s.dropToFile(); err != nil {
					log.Printf("failed to backup metrics to file; %v", err)
				}
			}
		}()
	}

	if storeInterval == 0 {
		s.isSyncBackup = true
	}

	go func() {
		<-ctx.Done()
		if err := s.dropToFile(); err != nil {
			log.Printf("failed to backup metrics on shutdown, %v", err)
		}
		if err := s.closeFile(); err != nil {
			log.Printf("failed to properly close file on shutdown, %v", err)
		}
	}()

	return nil
}

func (s *MemStore) closeFile() (err error) {
	return s.file.Close()
}

func (s *MemStore) dropToFile() (err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if err := s.file.Truncate(0); err != nil {
		return errors.Wrap(err, "can't truncate backup file before writing")
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return errors.Wrap(err, "can't reset I/O offset before writing")
	}

	if err = json.NewEncoder(s.file).Encode(&s.metrics); err != nil {
		return errors.Wrap(err, "can't drop memstore to file")
	}

	return nil
}

func (s *MemStore) restoreFromFile() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err = json.NewDecoder(s.file).Decode(&s.metrics); err != nil {
		return errors.Wrap(err, "can't restore memstore from file")
	}

	return nil
}
