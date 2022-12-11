package main

import (
	"bronya/internal/remote"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.HandleFunc("/{id}/{network:tcp[46]?}", remote.TCPHandler)
	})

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
