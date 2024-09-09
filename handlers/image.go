package handlers

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/tamimym/dicom-service/repositories"
)

func Image(repo repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceId := r.PathValue("instance")
		if instanceId == "" {
			http.Error(w, "No instance ID given", http.StatusBadRequest)
			return
		}

		dto, err := repo.Read(instanceId)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		if dto.ImagePath != "" {
			w.Header().Set("Content-Type", "image/png")
			http.ServeFile(w, r, dto.ImagePath)
			return
		}

		http.Error(w, fmt.Sprintf("No png file found for instance %s", instanceId), http.StatusNotFound)
	}
}
