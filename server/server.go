package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Rest struct {
	Hostname      string
	Port          string
	MetricHandler MetricHandler
}

func (s *Rest) Run(ctx context.Context) error {
	log.Printf("start server on %s:%s", s.Hostname, s.Port)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", s.Hostname, s.Port),
	}

	http.HandleFunc("/update/", s.MetricHandler.HandleMetricRecord(ctx))

	go func() {
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			log.Printf("failed to shutdown server, %v", err)
		}
	}()

	return srv.ListenAndServe()
}
