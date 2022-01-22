package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tyrylgin/collecter/agent"
	"github.com/tyrylgin/collecter/service/metric"
	"github.com/tyrylgin/collecter/storage/memstore"
)

var (
	pollInterval   = flag.Duration("poll_interval", time.Second*2, "Metric update interval")
	reportInterval = flag.Duration("report_interval", time.Second*10, "Metric push to server interval")

	hostname = flag.String("hostname", "http://127.0.0.1", "Hostname send metrics to")
	port     = flag.String("port", "8080", "Port for hostname send metrics to")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stopSignal := make(chan os.Signal, 1)
		signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-stopSignal
		cancel()
		os.Exit(0)
	}()

	memStore := memstore.NewStorage()
	agentSrv := agent.Service{
		ServerEndpoint: fmt.Sprintf("%s:%s/update/", *hostname, *port),
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		MetricSrv:      metric.NewProcessor(&memStore),
	}

	agentSrv.Start(ctx)
}
