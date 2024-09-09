package repositories

import (
	"errors"
	"fmt"
	"image/png"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/tamimym/dicom-service/models"
)

type FileRepository struct {
	uploadsDir string
	imageDir   string
}

func NewFileRepository(uploadsDir string, imageDir string) (Repository, error) {
	repo := &FileRepository{
		uploadsDir: uploadsDir,
		imageDir:   imageDir,
	}

	err := repo.initDir(uploadsDir)
	if err != nil {
		return nil, err
	}
	err = repo.initDir(imageDir)
	if err != nil {
		return nil, err
	}

	return repo, nil
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

	dto.ImagePath = repo.generateImage(dto)

	return nil
}

func (repo *FileRepository) Read(instanceId string) (*models.DicomDTO, error) {
	filename := fmt.Sprintf("%s.dcm", instanceId)

	dataset, err := dicom.ParseFile(filepath.Join(repo.uploadsDir, filename), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not read %s", filename))
		return nil, err
	}

	slog.Info("File successfully read", slog.String("filename", filename))

	imageFilePath := filepath.Join(repo.imageDir, fmt.Sprintf("%s.png", instanceId))
	_, err = os.Stat(imageFilePath)
	if err != nil {
		slog.Error(fmt.Sprintf("Instance %s has no png file", instanceId))
		return &models.DicomDTO{
			InstanceId: instanceId,
			Dataset:    &dataset,
		}, nil
	}

	return &models.DicomDTO{
		InstanceId: instanceId,
		Dataset:    &dataset,
		ImagePath:  imageFilePath,
	}, nil

}

func (repo *FileRepository) generateImage(dto *models.DicomDTO) string {
	pixelDataElement, err := dto.Dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return ""
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)

	fr := pixelDataInfo.Frames[0]
	i, err := fr.GetImage()
	if err != nil {
		slog.Error(fmt.Sprintf("Error while getting image: %v\n", err))
		return ""
	}

	filename := fmt.Sprintf("%s.png", dto.InstanceId)
	name := filepath.Join(repo.imageDir, filename)
	f, err := os.Create(name)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while creating file: %s", err.Error()))
		return ""
	}
	defer f.Close()

	err = png.Encode(f, i)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	slog.Info(fmt.Sprintf("Image %s written\n", name))

	return name
}

func (repo *FileRepository) initDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			slog.Info(fmt.Sprintf("Creating %s directory", dir))

			if err := os.Mkdir(dir, 0755); err != nil {
				slog.Error(fmt.Sprintf("Unable to create %s directory", dir))
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
