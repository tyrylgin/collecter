package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metricmock "github.com/tyrylgin/collecter/service/metric/mock"
)

func TestMetricHandler_HandleMetricRecord(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	mSrv := metricmock.NewMockProcessor(ctrl)
	mSrv.EXPECT().SetGauge(ctx, gomock.Any(), gomock.Any()).Return(nil)
	mSrv.EXPECT().IncreaseCounter(ctx, gomock.Any(), gomock.Any()).Return(nil)

	h, err := NewMetricHandler(mSrv)
	require.NoError(t, err)

	srv := Rest{
		Hostname:      "localhost",
		Port:          "5000",
		MetricHandler: *h,
	}
	testSrv := httptest.NewServer(srv.MetricHandler.HandleMetricRecord(ctx))
	defer testSrv.Close()

	resp, err := http.Post(testSrv.URL+"/update/1/m1/300", "text/plain", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(testSrv.URL+"/update/0/m1/300", "text/plain", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
