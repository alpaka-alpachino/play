package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/play/store"
)

func newRouter() *mux.Router {
	dbHandling := store.NewDBmetric()

	router := mux.NewRouter()
	router.Use(Middleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/{ServiceName}", dbHandling.MetricsCreate).Methods(http.MethodPost)
	api.HandleFunc("/successful", dbHandling.GetSuccessNumberFromAll).Methods(http.MethodGet)
	api.HandleFunc("/{ServiceName}/status", dbHandling.GetSuccessAndFailedForOne).Methods(http.MethodGet)
	api.HandleFunc("/{from}/{to}/status", dbHandling.HandledRequestsForDate).Methods(http.MethodGet)
	api.HandleFunc("/{ServiceName}", dbHandling.GetMetricsForService).Methods(http.MethodGet)

	return router

}
