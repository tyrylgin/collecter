package memstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyrylgin/collecter/model"
)

func TestMemStore_NewStorage(t *testing.T) {
	assert.Equal(t, &MemStore{metrics: model.MetricMap{}}, NewStorage())
}

func TestMemStore_GetAll(t *testing.T) {
	store := NewStorage()
	store.metrics = model.MetricMap{
		"m1": model.Counter{},
		"m2": model.Gauge{},
	}

	assert.Equal(t, store.metrics, store.GetAll())
}

func TestMemStore_GetByName(t *testing.T) {
	store := NewStorage()
	store.metrics = model.MetricMap{
		"m1": model.Counter{},
		"m2": model.Gauge{},
	}

	assert.Equal(t, store.metrics["m1"], store.Get("m1"))
	assert.Equalf(t, nil, store.Get("m3"), "must return nil for unlisted name")
}

func TestMemStore_Save(t *testing.T) {
	s := NewStorage()
	s.metrics = model.MetricMap{
		"m1": model.Counter{},
		"m2": model.Gauge{},
	}

	es := NewStorage()
	es.metrics = model.MetricMap{
		"m1": model.Counter{},
		"m2": model.Gauge{},
		"m3": model.Gauge{},
	}

	nc := model.Gauge{}
	_ = s.Save("m3", nc)
	assert.Equalf(t, es.metrics, s.metrics, "add new metric to memstore")

	_ = s.Save("m3", nc)
	assert.Equalf(t, es.metrics, s.metrics, "add metric with taken name do nothing")
}
