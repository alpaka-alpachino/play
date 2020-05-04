package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	address := flag.String("address", ":1433", "address of server")
	flag.Parse()
	dbHandling := newDBmetric()
	defer dbHandling.db.Close()

	router := mux.NewRouter()
	router.Use(Middleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/{ServiceName}", dbHandling.MetricsCreate).Methods(http.MethodPost)
	api.HandleFunc("/successful", dbHandling.GetSuccessNumberFromAll).Methods(http.MethodGet)
	api.HandleFunc("/{ServiceName}/status", dbHandling.GetSuccessAndFailedForOne).Methods(http.MethodGet)
	api.HandleFunc("/{from}/{to}/status", dbHandling.HandledRequestsForDate).Methods(http.MethodGet)
	api.HandleFunc("/{ServiceName}", dbHandling.GetMetricsForService).Methods(http.MethodGet)

	if err := http.ListenAndServe(*address, router); err != nil {
		log.Fatal(err.Error())
	}

}
