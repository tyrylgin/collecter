package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (s *Rest) Run(ctx context.Context, address string, port string) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
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
	router.Get("/", s.metricHandler.getAll)
	router.Post("/update/{metric_type}/{metric_name}/{metric_value}", s.metricHandler.processMetric)
	router.Get("/value/{metric_type}/{metric_name}", s.metricHandler.getMetricValue)
	return router
}
