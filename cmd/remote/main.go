package main

import (
	"bronya/internal/remote"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.HandleFunc("/{id}/{network:tcp[46]?}", remote.TCPHandler)
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	log.Printf("Listen Port: %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
