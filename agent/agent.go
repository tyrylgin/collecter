package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/tyrylgin/collecter/api"
	"github.com/tyrylgin/collecter/pkg/memstat"
	metricService "github.com/tyrylgin/collecter/service/metric"
)

type Service struct {
	ServerHost     string
	PollInterval   time.Duration
	ReportInterval time.Duration
	HashKey        string
	MetricSrv      metricService.Processor
}

func (srv Service) Run(ctx context.Context) {
	log.Printf(
		"start agent, collect metrics on %v, poll interval is %s, report interval is %s",
		srv.ServerHost, srv.PollInterval, srv.ReportInterval,
	)

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

	var reqBody []api.Metric

	for name, metric := range metrics {
		metricToSend := api.ModelToMetric(name, metric)
		if srv.HashKey != "" {
			metricToSend.CalcHash(srv.HashKey)
		}

		reqBody = append(reqBody, metricToSend)
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("failed to marshal metrics: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/updates/", srv.ServerHost), "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Printf("failed to send metrics value %s, %v", reqBodyJSON, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read response body; %v", err)
			return
		}
		log.Printf("failed to send metrics value; server respond: %v, %s", resp.StatusCode, body)
	}
}
