package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tyrylgin/collecter/internal/server"
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

	srv := server.Server{
		Hostname: *hostname,
		Port:     *port,
	}

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("can't start server, %v", err)
	}
}
