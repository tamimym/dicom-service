package handlers

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"

	"github.com/tamimym/dicom-service/models"
	"github.com/tamimym/dicom-service/repositories"
)

func QueryHeader(repo repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceId := r.PathValue("instance")
		if instanceId == "" {
			http.Error(w, "No instance ID given", http.StatusBadRequest)
			return
		}

		tagValue := r.URL.Query().Get("tag")
		tag, err := models.ParseTag(tagValue)
		if err != nil {
			http.Error(w, "Bad tag given", http.StatusBadRequest)
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

		element, err := dto.Dataset.FindElementByTag(tag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		w.Header().Set("content-type", "application/json")
		if err = json.NewEncoder(w).Encode(element); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
