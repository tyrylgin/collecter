package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/tyrylgin/collecter/api"
	"github.com/tyrylgin/collecter/pkg/memstat"
	metricService "github.com/tyrylgin/collecter/service/metric"
)

type Service struct {
	ServerEndpoint string
	PollInterval   time.Duration
	ReportInterval time.Duration
	MetricSrv      metricService.Processor
}

func (srv Service) Start(ctx context.Context) {
	log.Printf("start agent, collect metrics on %v", srv.ServerEndpoint)

	pollTicker := time.NewTicker(srv.PollInterval)
	reportTicker := time.NewTicker(srv.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			srv.SnapshotMetrics()
		case <-reportTicker.C:
			srv.SendMetrics()
		case <-ctx.Done():
			pollTicker.Stop()
			reportTicker.Stop()
		}
	}
}

type Agent interface {
	SnapshotMetrics(ctx context.Context)
	SendMetrics(ctx context.Context) error
}

func (srv Service) SnapshotMetrics() {
	rand.Seed(time.Now().UnixNano())
	if err := srv.MetricSrv.SetGauge("RandomValue", rand.Float64()); err != nil {
		log.Println(err)
	}

	if err := srv.MetricSrv.IncreaseCounter("PollCount", 1); err != nil {
		log.Println(err)
	}

	memStat := memstat.GetRuntimeMemstat()
	for name, value := range memStat {
		if err := srv.MetricSrv.SetGauge(name, value); err != nil {
			log.Println(err)
		}
	}
}

func (srv Service) SendMetrics() {
	metrics := srv.MetricSrv.GetAll()

	for name, metric := range metrics {
		metricToSend := api.ModelToMetric(name, metric)

		metricB, err := json.Marshal(metricToSend)
		if err != nil {
			log.Printf("failed to marshal metric %s, %v", name, err)
		}

		resp, err := http.Post(srv.ServerEndpoint, "application/json", bytes.NewBuffer(metricB))
		if err != nil {
			log.Printf("failed to send metric value %s, %v", metricB, err)
			continue
		}
		resp.Body.Close()
	}
}
