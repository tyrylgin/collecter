package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCounter_Type(t *testing.T) {
	c := &Counter{}
	r := c.Type()
	assert.Equal(t, MetricTypeCounter, r)
}
