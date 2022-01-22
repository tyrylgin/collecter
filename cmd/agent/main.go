package main

import (
	"context"
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

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
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
