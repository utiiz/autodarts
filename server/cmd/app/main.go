package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/utiiz/autodarts/internal/handlers"
)

func main() {
	router := chi.NewMux()

	router.Get("/", handlers.Make(handlers.Index))

	listenAddr := ":8080"
	slog.Info("HTTP Server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
