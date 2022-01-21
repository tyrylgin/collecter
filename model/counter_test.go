package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCounter_Count(t *testing.T) {
	tests := []struct {
		name  string
		count int64
		want  int64
	}{
		{
			name:  "common",
			count: 1,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultCounter{
				delta: tt.count,
			}

			r := c.GetDelta()
			assert.Equal(t, tt.want, r)
		})
	}
}

func TestDefaultCounter_Increase(t *testing.T) {
	tests := []struct {
		name  string
		count int64
		i     int64
		want  int64
	}{
		{
			name:  "common",
			count: 1,
			i:     3,
			want:  4,
		},
		{
			name:  "increase by zero",
			count: 1,
			i:     0,
			want:  1,
		},
		{
			name:  "increase by negative int",
			count: 1,
			i:     -1,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DefaultCounter{
				delta: tt.count,
			}
			c.IncreaseDelta(tt.i)
			assert.Equal(t, tt.want, c.delta)
		})
	}
}

func TestDefaultCounter_Type(t *testing.T) {
	c := &DefaultCounter{}
	r := c.Type()
	assert.Equal(t, MetricTypeCounter, r)
}

func TestNewCounter(t *testing.T) {
	c := &DefaultCounter{}
	r := NewCounter()
	assert.Equal(t, c, r)
}
