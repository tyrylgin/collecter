package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/tyrylgin/collecter/api"
	"github.com/tyrylgin/collecter/storage/memstore"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	IsRestore     bool          `env:"RESTORE" envDefault:"true"`
	StoreFile     string        `env:"TEMP_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
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

	store := memstore.NewStorage()
	if err := store.WithFileBackup(ctx, cfg.StoreFile, cfg.StoreInterval, cfg.IsRestore); err != nil {
		log.Fatalf("failed to init backup file for memstore; %v", err)
	}

	srv := api.Rest{}
	srv.WithStorage(&store)
	if err := srv.Run(ctx, cfg.Address); err != nil {
		log.Fatalf("can't start server, %v", err)
	}
}
