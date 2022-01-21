package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGauge_Value(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "common",
			value: 1,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DefaultGauge{
				value: tt.value,
			}

			r := g.GetValue()
			assert.Equal(t, tt.want, r)
		})
	}
}

func TestDefaultGauge_Set(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "common",
			value: 1,
			want:  500.3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DefaultGauge{
				value: tt.value,
			}
			g.SetValue(tt.want)
			assert.Equal(t, tt.want, g.value)
		})
	}
}

func TestDefaultGauge_Type(t *testing.T) {
	g := &DefaultGauge{}
	r := g.Type()
	assert.Equal(t, MetricTypeGauge, r)
}

func TestNewGauge(t *testing.T) {
	g := &DefaultGauge{}
	r := NewGauge()
	assert.Equal(t, g, r)
}
