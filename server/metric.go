package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/service/metric"
)

type MetricHandler struct {
	metricSrv metric.Processor
}

func NewMetricHandler(metricSrv metric.Processor) (*MetricHandler, error) {
	if metricSrv == nil {
		return nil, fmt.Errorf("metric.Processor: nil")
	}

	return &MetricHandler{metricSrv: metricSrv}, nil
}

func (h MetricHandler) HandleMetricRecord(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pathFragments := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")

		if len(pathFragments) < 3 {
			http.Error(rw, "wrong update metric url, must be /update/<type>/<name>/<value>", http.StatusBadRequest)
			return
		}

		metricName := pathFragments[1]
		metricType := pathFragments[0]

		switch model.MetricType(metricType) {
		case model.MetricTypeCounter:
			counterValue, err := strconv.ParseInt(pathFragments[2], 10, 64)
			if err != nil {
				http.Error(rw, "failed to parse counter value", http.StatusBadRequest)
				return
			}

			if err := h.metricSrv.IncreaseCounter(ctx, metricName, counterValue); err != nil {
				http.Error(rw, "failed to save counter value", http.StatusInternalServerError)
				return
			}
		case model.MetricTypeGauge:
			gaugeValue, err := strconv.ParseFloat(pathFragments[2], 64)
			if err != nil {
				http.Error(rw, "failed to parse gauge value", http.StatusBadRequest)
				return
			}

			if err := h.metricSrv.SetGauge(ctx, metricName, gaugeValue); err != nil {
				http.Error(rw, "failed to save gauge value", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(rw, "wrong metric type", http.StatusBadRequest)
			return
		}

		rw.Write(nil)
	}
}
