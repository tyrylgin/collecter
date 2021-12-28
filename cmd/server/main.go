package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tyrylgin/collecter/server"
	"github.com/tyrylgin/collecter/service/metric"
	"github.com/tyrylgin/collecter/storage/mem"
)

var (
	hostname = flag.String("hostname", "127.0.0.1", "Hostname to bind to")
	port     = flag.String("port", "8080", "Port to bind to")
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

	store := mem.NewStorage()
	metricHandler, err := server.NewMetricHandler(metric.NewProcessor(&store))
	if err != nil {
		log.Fatalf("failed metric handler init, %v", err)
	}

	srv := server.Rest{
		Hostname:      *hostname,
		Port:          *port,
		MetricHandler: *metricHandler,
	}
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("can't start server, %v", err)
	}
}
