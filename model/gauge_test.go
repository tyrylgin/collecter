package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGauge(t *testing.T) {
	g := Gauge{}
	r := g.Type()
	assert.Equal(t, MetricTypeGauge, r)
}
