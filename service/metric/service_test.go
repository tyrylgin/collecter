package metric

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrylgin/collecter/model"
	storagemock "github.com/tyrylgin/collecter/storage/mock"
)

func TestNewProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := storagemock.NewMockMetricStorer(ctrl)

	assert.Equal(t, &Service{metricStore: s}, NewProcessor(s))
}

func TestService_GetAll(t *testing.T) {
	ctx := context.TODO()

	expected := map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
	}

	ctrl := gomock.NewController(t)
	s := storagemock.NewMockMetricStorer(ctrl)
	s.EXPECT().GetAll(ctx).Return(expected)

	p := NewProcessor(s)

	assert.Equal(t, expected, p.GetAll(ctx))
}

func TestService_IncreaseCounter(t *testing.T) {
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	s := storagemock.NewMockMetricStorer(ctrl)
	p := NewProcessor(s)

	s.EXPECT().GetByName(ctx, gomock.Eq("m1")).Return(nil)
	s.EXPECT().Save(ctx, gomock.Eq("m1"), gomock.Any()).Return(nil)

	err := p.IncreaseCounter(ctx, "m1", 1)
	require.NoError(t, err)

	s.EXPECT().GetByName(ctx, gomock.Eq("m1")).Return(nil)
	s.EXPECT().Save(ctx, gomock.Eq("m1"), gomock.Any()).Return(errors.New("save error"))

	err = p.IncreaseCounter(ctx, "m1", 1)
	require.Error(t, err)
}

func TestService_SetGauge(t *testing.T) {
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	s := storagemock.NewMockMetricStorer(ctrl)
	p := NewProcessor(s)

	s.EXPECT().GetByName(ctx, gomock.Eq("m1")).Return(nil)
	s.EXPECT().Save(ctx, gomock.Eq("m1"), gomock.Any()).Return(nil)

	err := p.SetGauge(ctx, "m1", 1)
	require.NoError(t, err)

	s.EXPECT().GetByName(ctx, gomock.Eq("m1")).Return(nil)
	s.EXPECT().Save(ctx, gomock.Eq("m1"), gomock.Any()).Return(errors.New("save error"))

	err = p.SetGauge(ctx, "m1", 1)
	require.Error(t, err)
}
