package memstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyrylgin/collecter/model"
)

func TestMemStore_NewStorage(t *testing.T) {
	assert.Equal(t, MemStore{metrics: make(map[string]model.Metric)}, NewStorage())
}

func TestMemStore_GetAll(t *testing.T) {
	store := NewStorage()
	store.metrics = map[string]model.Metric{
		"m1": &model.DefaultCounter{},
		"m2": &model.DefaultGauge{},
	}

	assert.Equal(t, store.metrics, store.GetAll())
}

func TestMemStore_GetByName(t *testing.T) {
	store := NewStorage()
	store.metrics = map[string]model.Metric{
		"m1": &model.DefaultCounter{},
		"m2": &model.DefaultGauge{},
	}

	assert.Equal(t, store.metrics["m1"], store.Get("m1"))
	assert.Equalf(t, nil, store.Get("m3"), "must return nil for unlisted name")
}

func TestMemStore_Save(t *testing.T) {
	s := NewStorage()
	s.metrics = map[string]model.Metric{
		"m1": &model.DefaultCounter{},
		"m2": &model.DefaultGauge{},
	}

	es := NewStorage()
	es.metrics = map[string]model.Metric{
		"m1": &model.DefaultCounter{},
		"m2": &model.DefaultGauge{},
		"m3": &model.DefaultGauge{},
	}

	nc := &model.DefaultGauge{}
	_ = s.Save("m3", nc)
	assert.Equalf(t, es.metrics, s.metrics, "add new metric to memstore")

	_ = s.Save("m3", nc)
	assert.Equalf(t, es.metrics, s.metrics, "add metric with taken name do nothing")
}
