package repositories

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/suyashkumar/dicom"
	"github.com/tamimym/dicom-service/models"
)

type FileRepository struct {
	uploadsDir string
}

func NewFileRepository(uploadsDir string) (Repository, error) {
	if _, err := os.Stat(uploadsDir); err != nil {
		if os.IsNotExist(err) {
			slog.Info(fmt.Sprintf("Creating %s directory", uploadsDir))

			if err := os.Mkdir(uploadsDir, 0755); err != nil {
				slog.Error(fmt.Sprintf("Unable to create %s directory", uploadsDir))
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &FileRepository{
		uploadsDir: uploadsDir,
	}, nil
}

func (repo *FileRepository) Create(dto *models.DicomDTO) error {
	filename := fmt.Sprintf("%s.dcm", dto.InstanceId)

	dicomFile, err := os.Create(filepath.Join(repo.uploadsDir, filename))
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create file: %s in %s directory", filename, repo.uploadsDir))
		return err
	}
	defer dicomFile.Close()

	err = dicom.Write(dicomFile, *dto.Dataset)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not write to %s", filename))
		return err
	}

	slog.Info("File successfully written", slog.String("filename", filename))

	return nil
}
