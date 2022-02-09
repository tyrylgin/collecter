package api

import (
	"log"
	"net/http"

	"github.com/tyrylgin/collecter/storage"
	"github.com/tyrylgin/collecter/storage/psstore"
)

func PingDBHandler(store storage.MetricStorer) http.HandlerFunc {
	dbStore, ok := store.(*psstore.PsStore)

	return func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			log.Print("unsuccessful ping to database; database not initialized")
			http.Error(w, "database not initialized", http.StatusInternalServerError)
			return
		}

		if err := dbStore.Ping(); err != nil {
			log.Printf("unsuccessful ping to database: %v\n", err)
			http.Error(w, "unsuccessful ping to database", http.StatusInternalServerError)
			return
		}
	}
}
