package repositories

import "github.com/tamimym/dicom-service/models"

type Repository interface {
	Create(dto *models.DicomDTO) error
}
