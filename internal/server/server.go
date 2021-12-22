package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Hostname string
	Port     string
}

func (config *Server) Start(ctx context.Context) error {
	log.Printf("start server on %s:%s", config.Hostname, config.Port)

	http.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {})

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", config.Hostname, config.Port),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			log.Printf("failed to shutdown server, %v", err)
		}
	}()

	return srv.ListenAndServe()
}
