package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tamimym/dicom-service/handlers"
	"github.com/tamimym/dicom-service/repositories"
)

const (
	UPLOADS_DIR = "uploads"
)

func NewRouter(repo repositories.Repository) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /instance/{instance}", handlers.QueryHeader(repo))
	router.HandleFunc("POST /instance", handlers.Upload(repo))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	return router
}

func main() {
	repo, err := repositories.NewFileRepository(UPLOADS_DIR)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to setup file repository at: %s", UPLOADS_DIR))
		panic(err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: NewRouter(repo),
	}

	slog.Info("Starting server on port :8080")
	server.ListenAndServe()
}
