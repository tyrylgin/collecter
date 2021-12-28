package api

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrylgin/collecter/model"
	metricmock "github.com/tyrylgin/collecter/service/metric/mock"
)

func Test_metricHandler_getAll(t *testing.T) {
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().GetAll().Return(map[string]model.Metric{
		"m1": &model.DefaultGauge{},
		"m2": &model.DefaultCounter{},
	})
	mock.EXPECT().GetAll().Return(map[string]model.Metric{})
	rest := &Rest{metricHandler{metricService: mock}}

	ts := httptest.NewServer(rest.router())
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1 0\nm2 0\n", body)

	resp, body = testRequest(t, ts, "GET", "/")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "no metric registered yet\n", body)
}

func Test_metricHandler_getMetricValue(t *testing.T) {
	typeCounter := model.MetricTypeCounter
	typeGauge := model.MetricTypeGauge
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().Get("m1", &typeCounter).Return(&model.DefaultCounter{}, nil)
	mock.EXPECT().Get("m2", &typeCounter).Return(nil, errors.New("some error"))
	mock.EXPECT().Get("m3", &typeCounter).Return(nil, nil)
	mock.EXPECT().Get("m4", &typeGauge).Return(&model.DefaultGauge{}, nil)
	rest := &Rest{metricHandler{metricService: mock}}

	ts := httptest.NewServer(rest.router())
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/value/counter/m1")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "0", body)

	resp, body = testRequest(t, ts, "GET", "/value/unknown/m2")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
	assert.Equal(t, "wrong metric type\n", body)

	resp, body = testRequest(t, ts, "GET", "/value/counter/m2")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "failed to get metric\n", body)

	resp, body = testRequest(t, ts, "GET", "/value/counter/m3")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "metric not found\n", body)

	resp, body = testRequest(t, ts, "GET", "/value/gauge/m4")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "0", body)
}

func Test_metricHandler_processMetric(t *testing.T) {
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().IncreaseCounter("m1", gomock.Any()).Return(nil)
	mock.EXPECT().IncreaseCounter("m2", gomock.Any()).Return(errors.New("some error"))
	mock.EXPECT().SetGauge("m1", gomock.Any()).Return(nil)
	mock.EXPECT().SetGauge("m2", gomock.Any()).Return(errors.New("some error"))
	rest := &Rest{metricHandler{metricService: mock}}

	ts := httptest.NewServer(rest.router())
	defer ts.Close()

	resp, _ := testRequest(t, ts, "POST", "/update/counter/m1/1")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/unknown/m1/1")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/counter/m1/87yw")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, body := testRequest(t, ts, "POST", "/update/counter/m2/3")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "failed to set counter value\n", body)

	resp, _ = testRequest(t, ts, "POST", "/update/gauge/m1/1")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/unknown/m1/1")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/gauge/m1/87yw")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, body = testRequest(t, ts, "POST", "/update/gauge/m2/3")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "failed to set counter value\n", body)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("error during http request execution")
	}
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)

	require.NoError(t, err)

	return resp, string(respBody)
}
