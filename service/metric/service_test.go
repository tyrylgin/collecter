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
	s.EXPECT().Get("m1").AnyTimes().Return(&model.DefaultGauge{})
	p := NewProcessor(s)

	m, err := p.Get("m1", nil)
	require.NoError(t, err)
	assert.Equal(t, &model.DefaultGauge{}, m)

	typeCounter := model.MetricTypeCounter
	m, err = p.Get("m1", &typeCounter)
	require.Error(t, err)
	assert.Nil(t, m)
}

func TestService_GetAll(t *testing.T) {
	exp := map[string]model.Metric{
		"m1": &model.DefaultCounter{},
		"m2": &model.DefaultGauge{},
	}

	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().GetAll().Return(exp)
	p := NewProcessor(s)

	assert.Equal(t, exp, p.GetAll())
}

func TestService_IncreaseCounter(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().Get("m1").Return(&model.DefaultCounter{})
	s.EXPECT().Get("m2").Return(&model.DefaultGauge{})
	s.EXPECT().Get("m3").Return(nil)
	s.EXPECT().Save("m3", gomock.Any()).Return(nil)
	p := NewProcessor(s)

	err := p.IncreaseCounter("m1", 1)
	require.NoError(t, err)

	err = p.IncreaseCounter("m2", 1)
	require.Errorf(t, err, "error if trying increase gauge")

	err = p.IncreaseCounter("m3", 1)
	require.NoError(t, err, "save new counter if name not occupied")
}

func TestService_SetGauge(t *testing.T) {
	s := storagemock.NewMockMetricStorer(gomock.NewController(t))
	s.EXPECT().Get("m1").Return(&model.DefaultGauge{})
	s.EXPECT().Get("m2").Return(&model.DefaultCounter{})
	s.EXPECT().Get("m3").Return(nil)
	s.EXPECT().Save("m3", gomock.Any()).Return(nil)
	p := NewProcessor(s)

	err := p.SetGauge("m1", 1)
	require.NoError(t, err)

	err = p.SetGauge("m2", 1)
	require.Errorf(t, err, "error if trying set counter")

	err = p.SetGauge("m3", 1)
	require.NoError(t, err, "save new gauge if name not occupied")
}
