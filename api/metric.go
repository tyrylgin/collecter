package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tyrylgin/collecter/model"
	metricService "github.com/tyrylgin/collecter/service/metric"
)

type Metrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type metricHandler struct {
	metricService metricService.Processor
}

func (h *metricHandler) processMetricJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	var metric Metrics
	if err = json.Unmarshal(b, &metric); err != nil {
		log.Printf("failed to unmarshal metric json: %v", err)
		http.Error(w, "failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	if err := model.MetricType(metric.MType).Validate(); err != nil {
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}

	switch model.MetricType(metric.MType) {
	case model.MetricTypeCounter:
		if err = h.metricService.IncreaseCounter(metric.ID, metric.Delta); err != nil {
			http.Error(w, "failed to set counter value", http.StatusInternalServerError)
			return
		}
	case model.MetricTypeGauge:
		if err = h.metricService.SetGauge(metric.ID, metric.Value); err != nil {
			http.Error(w, "failed to set counter value", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(""))
}

func (h *metricHandler) processMetric(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "metric_name")
	mType := chi.URLParam(r, "metric_type")
	value := chi.URLParam(r, "metric_value")
	fmt.Printf("----- %s ----- %v -----\n", name, value)

	if err := model.MetricType(mType).Validate(); err != nil {
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}

	switch model.MetricType(mType) {
	case model.MetricTypeCounter:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "failed to parse counter value", http.StatusBadRequest)
			return
		}
		if err = h.metricService.IncreaseCounter(name, intValue); err != nil {
			http.Error(w, "failed to set counter value", http.StatusInternalServerError)
			return
		}
	case model.MetricTypeGauge:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "failed to parse counter value", http.StatusBadRequest)
			return
		}
		if err = h.metricService.SetGauge(name, floatValue); err != nil {
			http.Error(w, "failed to set counter value", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(""))
}

func (h *metricHandler) getMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "metric_name")
	mType := model.MetricType(chi.URLParam(r, "metric_type"))

	if err := mType.Validate(); err != nil {
		http.Error(w, "wrong metric type", http.StatusNotImplemented)
		return
	}

	metric, err := h.metricService.Get(name, &mType)
	if err != nil {
		http.Error(w, "failed to get metric", http.StatusBadRequest)
		return
	}
	if metric == nil {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	switch mType {
	case model.MetricTypeCounter:
		_, err = fmt.Fprintf(w, "%v", metric.(model.Counter).GetDelta())
	case model.MetricTypeGauge:
		_, err = fmt.Fprintf(w, "%v", metric.(model.Gauge).GetValue())
	}
	if err != nil {
		http.Error(w, "can't write response", http.StatusInternalServerError)
		return
	}
}

func (h *metricHandler) getAll(w http.ResponseWriter, r *http.Request) {
	metrics := h.metricService.GetAll()
	if len(metrics) == 0 {
		http.Error(w, "no metric registered yet", http.StatusNotFound)
		return
	}

	metricNames := make([]string, 0, len(metrics))
	for k := range metrics {
		metricNames = append(metricNames, k)
	}
	sort.Strings(metricNames)

	var metricsResp string
	for _, metricName := range metricNames {
		switch metrics[metricName].Type() {
		case model.MetricTypeCounter:
			metricsResp += fmt.Sprintf("%v %v\n", metricName, metrics[metricName].(model.Counter).GetDelta())
		case model.MetricTypeGauge:
			metricsResp += fmt.Sprintf("%v %v\n", metricName, metrics[metricName].(model.Gauge).GetValue())
		}
	}

	w.Write([]byte(metricsResp))
}
