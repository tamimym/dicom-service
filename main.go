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

func main() {
	repo, err := repositories.NewFileRepository(UPLOADS_DIR)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to setup file repository at: %s", UPLOADS_DIR))
		panic(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /studies/{study}/series/{series}/instance/{instance}/metadata", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	router.HandleFunc("POST /instance", handlers.Upload(repo))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	slog.Info("Starting server on port :8080")
	server.ListenAndServe()
}
