package api

import (
	"bytes"
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
	mock.EXPECT().GetAll().Return(model.MetricMap{
		"m1": model.Gauge{},
		"m2": model.Counter{},
	})
	mock.EXPECT().GetAll().Return(model.MetricMap{})
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
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().Get("m1", model.MetricTypeCounter).Return(model.Counter{}, nil)
	mock.EXPECT().Get("m2", model.MetricTypeCounter).Return(nil, errors.New("some error"))
	mock.EXPECT().Get("m3", model.MetricTypeCounter).Return(nil, nil)
	mock.EXPECT().Get("m4", model.MetricTypeGauge).Return(model.Gauge{}, nil)
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

func Test_metricHandler_processMetricJSON(t *testing.T) {
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().IncreaseCounter("m1", gomock.Any()).Return(nil)
	rest := &Rest{metricHandler{metricService: mock}}

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		req  []byte
		want want
	}{
		{
			name: "common",
			req:  []byte(`{"id":"m1", "type":"counter", "delta":1}`),
			want: want{
				code: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(tt.req))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(rest.metricHandler.processMetricJSON)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func Test_metricHandler_getMetricValueJSON(t *testing.T) {
	mock := metricmock.NewMockProcessor(gomock.NewController(t))
	mock.EXPECT().Get("m1", model.MetricTypeCounter).Return(model.Counter{}, nil)
	rest := &Rest{metricHandler{metricService: mock}}

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		req  []byte
		want want
	}{
		{
			name: "common",
			req:  []byte(`{"id":"m1", "type":"counter"}`),
			want: want{
				code:        200,
				response:    `{"id":"m1", "type":"counter", "delta":0}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(tt.req))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(rest.metricHandler.getMetricValueJSON)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.JSONEq(t, tt.want.response, w.Body.String())
		})
	}
}
