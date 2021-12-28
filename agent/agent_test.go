package agent

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrylgin/collecter/model"
	metricmock "github.com/tyrylgin/collecter/service/metric/mock"
)

func TestService_SendMetrics(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/gauge/m1/0", r.URL.Path)

		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))

	mSrv := metricmock.NewMockProcessor(ctrl)
	mSrv.EXPECT().GetAll(ctx).AnyTimes().Return(map[string]model.Metric{
		"m1": &model.DefaultGauge{},
	})

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	s := Service{
		ServerEndpoint: testSrv.URL + "/update/",
		MetricSrv:      mSrv,
	}
	s.SendMetrics(ctx)
	assert.Equalf(t, "", logBuf.String(), "no err log in stdout")

	s.ServerEndpoint = testSrv.URL
	s.SendMetrics(ctx)
	assert.Containsf(t, logBuf.String(), "failed", "must log err to stdout when failed on send")
}

func TestService_SnapshotMetrics(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	mSrv := metricmock.NewMockProcessor(ctrl)
	mSrv.EXPECT().SetGauge(ctx, gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mSrv.EXPECT().IncreaseCounter(ctx, gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	s := Service{MetricSrv: mSrv}
	s.SnapshotMetrics(ctx)

	assert.Empty(t, logBuf.String())
}
