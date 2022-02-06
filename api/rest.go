package api

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	metricService "github.com/tyrylgin/collecter/service/metric"
	"github.com/tyrylgin/collecter/storage"
)

type Rest struct {
	metricHandler metricHandler
}

func (s *Rest) WithStorage(store storage.MetricStorer) {
	s.metricHandler = metricHandler{
		metricService: metricService.NewProcessor(store),
	}
}

func (s *Rest) SetHashKey(key string) {
	s.metricHandler.hashKey = key
}

func (s *Rest) Run(ctx context.Context, address string) error {
	log.Printf("start server on %v", address)

	server := &http.Server{
		Addr:    address,
		Handler: s.router(),
	}

	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			log.Printf("failed to shutdown server, %v", err)
		}
	}()

	return server.ListenAndServe()
}

func (s *Rest) router() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Compress(5))

	router.Get("/", s.metricHandler.getAll)

	router.Post("/update/", s.metricHandler.processMetricJSON)
	router.Post("/update/{metric_type}/{metric_name}/{metric_value}", s.metricHandler.processMetric)

	router.Post("/value/", s.metricHandler.getMetricValueJSON)
	router.Get("/value/{metric_type}/{metric_name}", s.metricHandler.getMetricValue)

	return router
}
