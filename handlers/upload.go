package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/tamimym/dicom-service/models"
	"github.com/tamimym/dicom-service/repositories"
)

func Upload(repo repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 25<<20) // Setup 25MB size limit

		file, fileHeader, err := r.FormFile("instance")
		if err != nil {
			if errors.As(err, new(*http.MaxBytesError)) {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		defer file.Close()

		slog.Info("File uploaded", slog.Int64("filesize", fileHeader.Size))

		dto, err := models.NewDicomDTO(file, fileHeader.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = repo.Create(dto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("File stored successfully")

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(dto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
