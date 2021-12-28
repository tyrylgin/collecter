package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tyrylgin/collecter/api"
	"github.com/tyrylgin/collecter/storage/memstore"
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

	store := memstore.NewStorage()
	srv := api.Rest{}
	srv.WithStorage(&store)
	if err := srv.Run(ctx, *hostname, *port); err != nil {
		log.Fatalf("can't start server, %v", err)
	}
}
