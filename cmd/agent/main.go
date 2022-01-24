package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/tyrylgin/collecter/agent"
	"github.com/tyrylgin/collecter/service/metric"
	"github.com/tyrylgin/collecter/storage/memstore"
)

const (
	Address        string        = "127.0.0.1:8080"
	PollInterval   time.Duration = time.Second * 2
	ReportInterval time.Duration = time.Second * 10
)

type config struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stopSignal := make(chan os.Signal, 1)
		signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-stopSignal
		cancel()
		os.Exit(0)
	}()

	var cfg config

	flag.StringVar(&cfg.Address, "a", Address, "Hostname send metrics to")
	flag.DurationVar(&cfg.PollInterval, "p", PollInterval, "Metric update interval")
	flag.DurationVar(&cfg.ReportInterval, "r", ReportInterval, "Metric push to server interval")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("failed to parse env variables to config; %+v\n", err)
	}

	memStore := memstore.NewStorage()
	agentSrv := agent.Service{
		ServerHost:     fmt.Sprintf("http://%s", cfg.Address),
		PollInterval:   cfg.PollInterval,
		ReportInterval: cfg.ReportInterval,
		MetricSrv:      metric.NewProcessor(&memStore),
	}

	agentSrv.Run(ctx)
}
