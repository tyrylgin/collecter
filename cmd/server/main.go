package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/tyrylgin/collecter/api"
	"github.com/tyrylgin/collecter/storage/memstore"
	"github.com/tyrylgin/collecter/storage/psstore"
)

const (
	Address       string        = "127.0.0.1:8080"
	IsRestore     bool          = true
	StoreFile     string        = "/tmp/devops-metrics-db.json"
	StoreInterval time.Duration = time.Second * 300
)

type config struct {
	Address       string        `env:"ADDRESS"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
	IsRestore     bool          `env:"RESTORE"`
	StoreFile     string        `env:"STORE_FILE"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	SecretKey     string        `env:"KEY"`
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
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database source name")
	flag.BoolVar(&cfg.IsRestore, "r", IsRestore, "Is restore from backup file")
	flag.StringVar(&cfg.StoreFile, "f", StoreFile, "Backup file path")
	flag.DurationVar(&cfg.StoreInterval, "i", StoreInterval, "Backup to file interval")
	flag.StringVar(&cfg.SecretKey, "k", "", "Secret key to sign data")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("failed to parse env variables to config; %+v\n", err)
	}

	srv := api.Rest{}
	if cfg.DatabaseDSN != "" {
		dbStore, err := psstore.Init(cfg.DatabaseDSN)
		if err != nil {
			log.Fatalf("unable to connect to database: %v\n", err)
		}

		srv.WithStorage(dbStore)
	} else {
		memStore := memstore.NewStorage()
		if err := memStore.WithFileBackup(ctx, cfg.StoreFile, cfg.StoreInterval, cfg.IsRestore); err != nil {
			log.Fatalf("failed to init backup file for memstore; %v", err)
		}

		srv.WithStorage(memStore)
	}

	srv.SetHashKey(cfg.SecretKey)
	if err := srv.Run(ctx, cfg.Address); err != nil {
		log.Fatalf("can't start server, %v", err)
	}
}
