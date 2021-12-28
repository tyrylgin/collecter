package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Rest struct {
	Hostname      string
	Port          string
	MetricHandler MetricHandler
}

func (s *Rest) Run(ctx context.Context) error {
	log.Printf("start server on %s:%s", s.Hostname, s.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.Hostname, s.Port),
		Handler: s.router(ctx),
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

func (s *Rest) router(ctx context.Context) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Get("/", s.MetricHandler.GetAll(ctx))
	router.Route("/update", func(r chi.Router) {
		r.Post("/{metric_type}/{metric_name}/{metric_value}", s.MetricHandler.HandleMetricRecord(ctx))
	})
	router.Route("/value", func(r chi.Router) {
		r.Get("/{metric_type}/{metric_name}", s.MetricHandler.GetMetricValue(ctx))
	})

	return router
}
