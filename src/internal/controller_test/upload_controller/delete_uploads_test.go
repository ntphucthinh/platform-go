package upload_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/usecases/uploads"
)

func TestDeleteUploads(t *testing.T) {
	client := &fakeAssetClient{}
	response := performDelete(t, newUploadRouter(t, client), `[{"path":"travel_location/550e8400-e29b-41d4-a716-446655440000.mp4","resourceType":"video"}]`)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())

	require.Len(t, client.deletes, 1)
	assert.Equal(t, "travel_location/550e8400-e29b-41d4-a716-446655440000", client.deletes[0].PublicID)
	assert.Equal(t, uploads.ResourceVideo, client.deletes[0].ResourceType)

	var body api.DeleteUploadsResponse
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Len(t, body.Data, 1)
	assert.True(t, body.Data[0].Success)
}

func TestDeleteUploadsDoesNotDeleteInvalidPath(t *testing.T) {
	client := &fakeAssetClient{}
	response := performDelete(t, newUploadRouter(t, client), `[{"path":"travel_location/not-a-uuid.mp4","resourceType":"video"}]`)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	assert.Empty(t, client.deletes)
	assert.Contains(t, response.Body.String(), "invalid delete request")
}

func performDelete(t *testing.T, router http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, "/uploads", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
