package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tamimym/dicom-service/handlers"
)

func TestQueryHeader(t *testing.T) {
	repo, dir1, dir2 := setup()
	defer teardown(dir1, dir2)

	t.Run("it returns bad request status if no instance id given", func(t *testing.T) {
		queryHeaderHandler := handlers.QueryHeader(repo)

		assert.HTTPStatusCode(t, queryHeaderHandler, http.MethodGet, "/instance/", nil, http.StatusBadRequest)
	})

	t.Run("it returns bad request status if no tag given", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/instance/test-instance-id", nil)
		request.SetPathValue("instance", "test-instance-id")
		response := httptest.NewRecorder()

		queryHeaderHandler := handlers.QueryHeader(repo)
		queryHeaderHandler(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})

	t.Run("it returns not found status if instance does not exist", func(t *testing.T) {
		queryValues := url.Values{}
		queryValues.Add("tag", "(0008,0021)")
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/instance/does-not-exist-id?%s", queryValues.Encode()), nil)
		request.SetPathValue("instance", "does-not-exist-id")
		response := httptest.NewRecorder()

		queryHeaderHandler := handlers.QueryHeader(repo)
		queryHeaderHandler(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	})

	t.Run("it returns header attribute based on tag given", func(t *testing.T) {
		// Copy a test file into the file repository path
		instanceId := "1.2.826.0.1.3680043.2.1074.8138928452617025399165543974931135073"
		input, _ := os.ReadFile("../test_data/IM000002")
		err := os.WriteFile(filepath.Join(dir1, fmt.Sprintf("%s.dcm", instanceId)), input, 0644)
		assert.Nil(t, err)

		// Setup the request
		queryValues := url.Values{}
		queryValues.Add("tag", "(0008,0090)")
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/instance/%s?%s", instanceId, queryValues.Encode()), nil)
		request.SetPathValue("instance", instanceId)
		response := httptest.NewRecorder()

		queryHeaderHandler := handlers.QueryHeader(repo)
		queryHeaderHandler(response, request)

		responseBody, err := io.ReadAll(response.Result().Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
		assert.Equal(t, "{\"tag\":{\"Group\":8,\"Element\":144},\"VR\":0,\"rawVR\":\"PN\",\"valueLength\":16,\"value\":[\"BROOKE DIX^DPM^\"]}\n", string(responseBody[:]))
	})
}
