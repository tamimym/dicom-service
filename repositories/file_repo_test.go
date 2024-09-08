package repositories_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamimym/dicom-service/models"
	"github.com/tamimym/dicom-service/repositories"
)

func setup() (repositories.Repository, string, *models.DicomDTO) {
	dir, _ := os.MkdirTemp("", "")
	repo, _ := repositories.NewFileRepository(dir)

	testFile, _ := os.Open("../test_data/IM000020")
	defer testFile.Close()

	fileInfo, _ := os.Stat("../test_data/IM000020")

	dto, _ := models.NewDicomDTO(testFile, fileInfo.Size())

	return repo, dir, dto
}

func teardown(dir string) {
	os.RemoveAll(dir)
}

func TestFileRepository(t *testing.T) {
	repo, dir, dto := setup()
	defer teardown(dir)

	t.Run("it writes a dicom file to the directory with instance id as filename", func(t *testing.T) {
		err := repo.Create(dto)

		assert.Nil(t, err)
		assert.FileExists(t, filepath.Join(dir, fmt.Sprintf("%s.dcm", dto.InstanceId)))
	})
}
