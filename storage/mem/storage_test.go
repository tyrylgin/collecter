package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyrylgin/collecter/model"
)

func TestNewStorage(t *testing.T) {
	assert.Equal(t, MemStore{metrics: make(map[string]model.Metric)}, NewStorage())
}
