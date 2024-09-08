package models_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamimym/dicom-service/models"
)

func TestDicomDTO(t *testing.T) {
	t.Run("it parses a dicom file and returns the DTO", func(t *testing.T) {
		testFile, err := os.Open("../test_data/IM000012")
		assert.Nil(t, err)
		defer testFile.Close()

		fileInfo, err := os.Stat("../test_data/IM000012")
		assert.Nil(t, err)

		dto, err := models.NewDicomDTO(testFile, fileInfo.Size())

		assert.Nil(t, err)
		assert.Equal(t, "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000436", dto.InstanceId)
	})

	t.Run("it fails to parse a non-dicom file", func(t *testing.T) {
		testFile, err := os.Open("../test_data/example.dat")
		assert.Nil(t, err)
		defer testFile.Close()

		fileInfo, err := os.Stat("../test_data/example.dat")
		assert.Nil(t, err)

		_, err = models.NewDicomDTO(testFile, fileInfo.Size())

		assert.EqualError(t, err, "SOP Instance UID not found")
	})
}
