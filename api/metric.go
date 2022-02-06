package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tyrylgin/collecter/model"
	metricService "github.com/tyrylgin/collecter/service/metric"
)

type Metric struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func (m *Metric) CalcHash(hashKey string) {
	source := fmt.Sprintf("%s:%s", m.ID, m.Type)

	if m.Delta != nil {
		source = fmt.Sprintf("%s:%d", source, m.Delta)
	}

	if m.Value != nil {
		source = fmt.Sprintf("%s:%v", source, m.Value)
	}

	m.Hash = GetHash(source, hashKey)
}

func EqualHash(m Metric, hashKey string) bool {
	originalHash := m.Hash
	m.CalcHash(hashKey)

	return originalHash == m.Hash
}

func ModelToMetric(name string, m model.Metric) Metric {
	metric := Metric{
		ID:   name,
		Type: string(m.Type()),
	}

	switch m.Type() {
	case model.MetricTypeCounter:
		delta := m.(model.Counter).Delta
		metric.Delta = &delta
	case model.MetricTypeGauge:
		value := m.(model.Gauge).Value
		metric.Value = &value
	}

	return metric
}

func GetHash(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type metricHandler struct {
	hashKey       string
	metricService metricService.Processor
}

func (h *metricHandler) processMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		log.Printf("failed to unmarshal metric json: %v", err)
		http.Error(w, "failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	if h.hashKey != "" && EqualHash(metric, h.hashKey) {
		http.Error(w, "hashes not equal", http.StatusBadRequest)
		return
	}

	if err := model.MetricType(metric.Type).Validate(); err != nil {
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}

	switch model.MetricType(metric.Type) {
	case model.MetricTypeCounter:
		if err := h.metricService.IncreaseCounter(metric.ID, *metric.Delta); err != nil {
			http.Error(w, "failed to set counter value", http.StatusInternalServerError)
			return
		}
	case model.MetricTypeGauge:
		if err := h.metricService.SetGauge(metric.ID, *metric.Value); err != nil {
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

func (h *metricHandler) getMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	var reqMetric Metric
	if err := json.NewDecoder(r.Body).Decode(&reqMetric); err != nil {
		log.Printf("failed to unmarshal metric json: %v", err)
		http.Error(w, "failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	if err := model.MetricType(reqMetric.Type).Validate(); err != nil {
		http.Error(w, "wrong metric type", http.StatusNotImplemented)
		return
	}

	metric, err := h.metricService.Get(reqMetric.ID, model.MetricType(reqMetric.Type))
	if err != nil {
		http.Error(w, "failed to get metric", http.StatusBadRequest)
		return
	}
	if metric == nil {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	respMetric := ModelToMetric(reqMetric.ID, metric)

	if h.hashKey != "" {
		respMetric.CalcHash(h.hashKey)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(respMetric); err != nil {
		http.Error(w, "failed to marshal metric json", http.StatusInternalServerError)
		return
	}
}

func (h *metricHandler) getMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "metric_name")
	mType := model.MetricType(chi.URLParam(r, "metric_type"))

	w.Header().Set("Content-Type", "text/html")

	if err := mType.Validate(); err != nil {
		http.Error(w, "wrong metric type", http.StatusNotImplemented)
		return
	}

	metric, err := h.metricService.Get(name, mType)
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
		_, err = fmt.Fprintf(w, "%v", metric.(model.Counter).Delta)
	case model.MetricTypeGauge:
		_, err = fmt.Fprintf(w, "%v", metric.(model.Gauge).Value)
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
		metric := metrics[metricName]
		switch metric.Type() {
		case model.MetricTypeCounter:
			metricsResp += fmt.Sprintf("%v %v\n", metricName, metric.(model.Counter).Delta)
		case model.MetricTypeGauge:
			metricsResp += fmt.Sprintf("%v %v\n", metricName, metric.(model.Gauge).Value)
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(metricsResp))
}
