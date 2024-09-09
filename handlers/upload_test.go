package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamimym/dicom-service/handlers"
	"github.com/tamimym/dicom-service/models"
	"github.com/tamimym/dicom-service/repositories"
)

func setup() (repositories.Repository, string, string) {
	upload_dir, _ := os.MkdirTemp("", "")
	image_dir, _ := os.MkdirTemp("", "")
	repo, _ := repositories.NewFileRepository(upload_dir, image_dir)

	return repo, upload_dir, image_dir
}

func teardown(upload_dir string, image_dir string) {
	os.RemoveAll(upload_dir)
	os.RemoveAll(image_dir)
}

func TestUpload(t *testing.T) {
	repo, dir1, dir2 := setup()
	defer teardown(dir1, dir2)

	t.Run("it returns bad request status if no instance file is uploaded", func(t *testing.T) {
		uploadHandler := handlers.Upload(repo)

		assert.HTTPStatusCode(t, uploadHandler, http.MethodPost, "/instance", nil, http.StatusBadRequest)
	})

	t.Run("it returns entity too large status if file is too big", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("instance", "toobig.dat")
		buf, _ := os.Open("../test_data/toobig.dat")
		defer buf.Close()
		_, err := io.Copy(fw, buf)
		assert.Nil(t, err)
		writer.Close()

		request, _ := http.NewRequest(http.MethodPost, "/instance", body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		uploadHandler := handlers.Upload(repo)
		uploadHandler(response, request)

		assert.Equal(t, http.StatusRequestEntityTooLarge, response.Result().StatusCode)
	})

	t.Run("it returns response with instance ID if file upload successful", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("instance", "IM000002")
		buf, _ := os.Open("../test_data/IM000002")
		defer buf.Close()
		_, err := io.Copy(fw, buf)
		assert.Nil(t, err)
		writer.Close()

		request, _ := http.NewRequest(http.MethodPost, "/instance", body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		uploadHandler := handlers.Upload(repo)
		uploadHandler(response, request)

		responseBody, err := io.ReadAll(response.Result().Body)
		assert.Nil(t, err)

		var dicom models.DicomDTO
		err = json.Unmarshal(responseBody, &dicom)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
		assert.Equal(t, "1.2.826.0.1.3680043.2.1074.8138928452617025399165543974931135073", dicom.InstanceId)
	})
}
