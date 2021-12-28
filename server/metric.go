package server

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tyrylgin/collecter/model"
	metricService "github.com/tyrylgin/collecter/service/metric"
)

type MetricHandler struct {
	metricSrv metricService.Processor
}

func NewMetricHandler(metricSrv metricService.Processor) (*MetricHandler, error) {
	if metricSrv == nil {
		return nil, fmt.Errorf("metricService.Processor: nil")
	}

	return &MetricHandler{metricSrv: metricSrv}, nil
}

func (h MetricHandler) HandleMetricRecord(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metric_name")
		metricType := chi.URLParam(r, "metric_type")
		metricValue := chi.URLParam(r, "metric_value")

		switch model.MetricType(metricType) {
		case model.MetricTypeCounter:
			counterValue, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(rw, "failed to parse counter value", http.StatusBadRequest)
				return
			}

			if err := h.metricSrv.IncreaseCounter(ctx, metricName, counterValue); err != nil {
				http.Error(rw, "failed to save counter value", http.StatusInternalServerError)
				return
			}
		case model.MetricTypeGauge:
			gaugeValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(rw, "failed to parse gauge value", http.StatusBadRequest)
				return
			}

			if err := h.metricSrv.SetGauge(ctx, metricName, gaugeValue); err != nil {
				http.Error(rw, "failed to save gauge value", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(rw, "wrong metric type", http.StatusNotImplemented)
			return
		}

		rw.Write(nil)
	}
}

func (h MetricHandler) GetMetricValue(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metric_name")
		metricType := model.MetricType(chi.URLParam(r, "metric_type"))

		if err := metricType.Validate(); err != nil {
			http.Error(rw, "wrong metric type", http.StatusNotImplemented)
			return
		}

		metric, err := h.metricSrv.Get(ctx, metricName, &metricType)
		if err != nil {
			http.Error(rw, "failed to get metric", http.StatusBadRequest)
			return
		}

		if metric == nil {
			http.Error(rw, "metric not found", http.StatusNotFound)
			return
		}

		switch metricType {
		case model.MetricTypeCounter:
			_, err = fmt.Fprintf(rw, "%v", metric.(model.Counter).Count())
		case model.MetricTypeGauge:
			_, err = fmt.Fprintf(rw, "%v", metric.(model.Gauge).Value())
		}
		if err != nil {
			http.Error(rw, "can't write response", http.StatusInternalServerError)
			return
		}
	}
}

func (h MetricHandler) GetAll(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var err error
		var metricsResp string

		metrics := h.metricSrv.GetAll(ctx)
		if len(metrics) == 0 {
			http.Error(rw, "no metric registered yet", http.StatusNotFound)
			return
		}

		metricNames := make([]string, 0, len(metrics))
		for k := range metrics {
			metricNames = append(metricNames, k)
		}
		sort.Strings(metricNames)

		for _, metricName := range metricNames {
			switch metrics[metricName].Type() {
			case model.MetricTypeCounter:
				metricsResp += fmt.Sprintf("%v %v\n", metricName, metrics[metricName].(model.Counter).Count())
			case model.MetricTypeGauge:
				metricsResp += fmt.Sprintf("%v %v\n", metricName, metrics[metricName].(model.Gauge).Value())
			}

			if err != nil {
				http.Error(rw, "can't write response", http.StatusInternalServerError)
				return
			}
		}

		rw.Write([]byte(metricsResp))
	}
}
