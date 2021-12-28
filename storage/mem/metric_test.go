package mem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrylgin/collecter/model"
)

func TestMemStore_GetAll(t *testing.T) {
	ctx := context.TODO()
	store := NewStorage()
	store.metrics = map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
	}

	assert.Equal(t, store.metrics, store.GetAll(ctx))
}

func TestMemStore_GetByName(t *testing.T) {
	ctx := context.TODO()
	store := NewStorage()
	store.metrics = map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
	}

	assert.Equal(t, store.metrics["m1"], store.GetByName(ctx, "m1"))
	assert.Equalf(t, nil, store.GetByName(ctx, "m3"), "must return nil for unlisted name")
}

func TestMemStore_Save(t *testing.T) {
	ctx := context.TODO()
	store := NewStorage()
	store.metrics = map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
	}

	expectedStore := NewStorage()
	expectedStore.metrics = map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
		"m3": &model.DefaultCounter{},
	}

	newCounter := model.NewCounter()
	err := store.Save(ctx, "m3", newCounter.(model.Metric))
	require.NoError(t, err)
	assert.Equal(t, expectedStore.metrics, store.metrics)
}
