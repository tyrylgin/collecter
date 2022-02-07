package agent

import (
	"bytes"
	"io"
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
	ctrl := gomock.NewController(t)

	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/updates/", r.URL.Path)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("failed to read request body")
		}

		assert.Equal(t, `{"metrics":[{"id":"m1","type":"gauge","value":0},{"id":"m2","type":"counter","delta":0}]}`, string(b))

		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(nil)
		require.NoError(t, err)
	}))

	mSrv := metricmock.NewMockProcessor(ctrl)
	mSrv.EXPECT().GetAll().AnyTimes().Return(model.MetricMap{
		"m1": model.Gauge{},
		"m2": model.Counter{},
	})

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	s := Service{
		ServerHost: testSrv.URL,
		MetricSrv:  mSrv,
	}
	s.SendMetrics()
	assert.Equalf(t, "", logBuf.String(), "no err log in stdout")
}

func TestService_SnapshotMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)

	mSrv := metricmock.NewMockProcessor(ctrl)
	mSrv.EXPECT().SetGauge(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mSrv.EXPECT().IncreaseCounter(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	s := Service{MetricSrv: mSrv}
	s.SnapshotMetrics()

	assert.Empty(t, logBuf.String())
}
