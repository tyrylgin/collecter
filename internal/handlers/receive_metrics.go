package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tyrylgin/collecter/internal/metrics"
	"github.com/tyrylgin/collecter/internal/storage"
)

func ReceiveMetricsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pathFragments := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")

		if len(pathFragments) < 3 {
			http.Error(rw, "wrong update metric url, must be /update/<type>/<name>/<value>", http.StatusBadRequest)
			return
		}

		metricName := pathFragments[1]
		metricType, err := strconv.Atoi(pathFragments[0])
		if err != nil {
			http.Error(rw, "failed to parse metric type", http.StatusBadRequest)
			return
		}

		switch metrics.MetricType(metricType) {
		case metrics.MetricTypeCounter:
			counterValue, err := strconv.ParseInt(pathFragments[2], 10, 64)
			if err != nil {
				http.Error(rw, "failed to parse counter value", http.StatusBadRequest)
				return
			}

			storage.StoreCounter(metricName, counterValue)
		case metrics.MetricTypeGauge:
			gaugeValue, err := strconv.ParseFloat(pathFragments[2], 64)
			if err != nil {
				http.Error(rw, "failed to parse gauge value", http.StatusBadRequest)
				return
			}

			storage.StoreGauge(metricName, gaugeValue)
		default:
			http.Error(rw, "wrong metric type", http.StatusBadRequest)
			return
		}

		rw.Write(nil)
	}
}
