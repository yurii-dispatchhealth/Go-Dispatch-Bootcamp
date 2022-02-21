package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type readerController interface {
	Fetch(w http.ResponseWriter, r *http.Request)
	FetchConcurrently(w http.ResponseWriter, r *http.Request)
	FetchById(w http.ResponseWriter, r *http.Request)
	Feed(w http.ResponseWriter, r *http.Request)
	UpdateUsersFromFeed(w http.ResponseWriter, r *http.Request)
}

func Setup(c readerController) *mux.Router {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/fetch", c.Fetch).Methods(http.MethodGet).Name("Fetch")
	v1.HandleFunc("/fetch-concurrently", c.FetchConcurrently).Methods(http.MethodGet).Name("Fetch")
	v1.HandleFunc("/fetch/{id}", c.FetchById).Methods(http.MethodGet).Name("FetchById")
	v1.HandleFunc("/feed", c.Feed).Methods(http.MethodGet).Name("Feed")
	v1.HandleFunc("/run-update-users-from-feed", c.UpdateUsersFromFeed).Methods(http.MethodGet).Name("UpdateUsersFromFeed")

	return r
}
