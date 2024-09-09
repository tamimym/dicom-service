package repositories_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamimym/dicom-service/models"
	"github.com/tamimym/dicom-service/repositories"
)

func setup() (repositories.Repository, string, string, *models.DicomDTO) {
	dir, _ := os.MkdirTemp("", "")
	dir2, _ := os.MkdirTemp("", "")
	repo, _ := repositories.NewFileRepository(dir, dir2)

	testFile, _ := os.Open("../test_data/IM000020")
	defer testFile.Close()

	fileInfo, _ := os.Stat("../test_data/IM000020")

	dto, _ := models.NewDicomDTO(testFile, fileInfo.Size())

	return repo, dir, dir2, dto
}

func teardown(dir1 string, dir2 string) {
	os.RemoveAll(dir1)
	os.RemoveAll(dir2)
}

func TestFileRepository(t *testing.T) {
	repo, dir1, dir2, dto := setup()
	defer teardown(dir1, dir2)

	t.Run("it writes a dicom file to the directory with instance id as filename", func(t *testing.T) {
		err := repo.Create(dto)

		assert.Nil(t, err)
		assert.FileExists(t, filepath.Join(dir1, fmt.Sprintf("%s.dcm", dto.InstanceId)))
	})

	t.Run("it returns error if it tries to read an instance id that does not exist", func(t *testing.T) {
		d, err := repo.Read("abcd")

		assert.Nil(t, d)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("it returns DTO if read is successful", func(t *testing.T) {
		err := repo.Create(dto)
		assert.Nil(t, err)

		d, err := repo.Read(dto.InstanceId)

		assert.Nil(t, err)
		assert.Equal(t, dto.InstanceId, d.InstanceId)
		assert.NotNil(t, d.Dataset)
	})
}
