package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tyrylgin/collecter/internal/metrics"
)

var (
	pollInterval   = flag.Duration("poll_interval", time.Second*2, "Metric update interval")
	reportInterval = flag.Duration("report_interval ", time.Second*10, "Metric push to server interval")

	hostname = flag.String("hostname", "http://127.0.0.1", "Hostname to bind to")
	port     = flag.String("port", "8080", "Port to bind to")
)

func main() {
	flag.Parse()

	url := fmt.Sprintf("%s:%s/update/", *hostname, *port)

	pollTicker := time.NewTicker(*pollInterval)
	reportTicker := time.NewTicker(*reportInterval)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Create and register custom metrics
	pollCounter := metrics.NewCounter()
	randomValue := metrics.NewGauge()

	metrics.Register("PollCount", pollCounter)
	metrics.Register("RandomValue", randomValue)

	for {
		select {
		case <-pollTicker.C:
			pollCounter.Increase(1)

			rand.Seed(time.Now().UnixNano())
			randomValue.Set(rand.Float64())

			metrics.SnapshotRuntimeMetrics()
		case <-reportTicker.C:
			metrics.Each(func(metricName string, metric metrics.Metric) {
				metricLogString := fmt.Sprintf("%d/%s/", metric.Type(), metricName)

				switch metric.Type() {
				case metrics.MetricTypeCounter:
					metricLogString += fmt.Sprint(metric.(metrics.Counter).Count())
				case metrics.MetricTypeGauge:
					metricLogString += fmt.Sprint(metric.(metrics.Gauge).Value())
				}

				endpoint := url + metricLogString
				resp, err := http.Post(endpoint, "text/plain", nil)
				if err != nil {
					log.Println(err)
					return
				}
				err = resp.Body.Close()
				if err != nil {
					log.Println(err)
					return
				}
			})
		case <-stopSignal:
			os.Exit(0)
		}
	}
}
