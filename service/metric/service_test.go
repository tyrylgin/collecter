package metric

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrylgin/collecter/model"
	storagemock "github.com/tyrylgin/collecter/storage/mock"
)

func TestNewProcessor(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	assert.Equal(t, &Service{store: s}, NewProcessor(s))
}

func TestService_Get(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().Get("m1").AnyTimes().Return(&model.Gauge{})
	p := NewProcessor(s)

	m, err := p.Get("m1", model.MetricTypeGauge)
	require.NoError(t, err)
	assert.Equal(t, &model.Gauge{}, m)

	m, err = p.Get("m1", model.MetricTypeCounter)
	require.Error(t, err)
	assert.Nil(t, m)
}

func TestService_GetAll(t *testing.T) {
	exp := map[string]model.Metric{
		"m1": &model.Counter{},
		"m2": &model.Gauge{},
	}

	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().GetAll().Return(exp)
	p := NewProcessor(s)

	assert.Equal(t, exp, p.GetAll())
}

func TestService_IncreaseCounter(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().Get("m1").Return(&model.Counter{})
	s.EXPECT().Save("m1", gomock.Any()).Return(nil)
	s.EXPECT().Get("m3").Return(nil)
	s.EXPECT().Save("m3", gomock.Any()).Return(nil)
	p := NewProcessor(s)

	err := p.IncreaseCounter("m1", 1)
	require.NoError(t, err)

	err = p.IncreaseCounter("m3", 1)
	require.NoError(t, err, "save new counter if name not occupied")
}

func TestService_SetGauge(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().Get("m1").Return(&model.Gauge{})
	s.EXPECT().Save("m1", gomock.Any()).Return(nil)
	s.EXPECT().Get("m3").Return(nil)
	s.EXPECT().Save("m3", gomock.Any()).Return(nil)
	p := NewProcessor(s)

	err := p.SetGauge("m1", 1)
	require.NoError(t, err)

	err = p.SetGauge("m3", 1)
	require.NoError(t, err, "save new gauge if name not occupied")
}
