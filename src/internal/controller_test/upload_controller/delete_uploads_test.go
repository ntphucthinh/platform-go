package upload_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "platform-go/generated/api"
	"platform-go/internal/usecases/uploads"
)

func TestDeleteUploads(t *testing.T) {
	client := &fakeAssetClient{}
	response := performDelete(t, newUploadRouter(t, client), `{"data":[{"path":"travel/550e8400-e29b-41d4-a716-446655440000.mp4","resourceType":"video"}]}`)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if len(client.deletes) != 1 || client.deletes[0].PublicID != "travel/550e8400-e29b-41d4-a716-446655440000" || client.deletes[0].ResourceType != uploads.ResourceVideo {
		t.Fatalf("delete request = %#v", client.deletes)
	}

	var body api.DeleteUploadsResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || !body.Data[0].Success {
		t.Fatalf("response = %#v", body)
	}
}

func TestDeleteUploadsDoesNotDeleteInvalidPath(t *testing.T) {
	client := &fakeAssetClient{}
	response := performDelete(t, newUploadRouter(t, client), `{"data":[{"path":"travel/not-a-uuid.mp4","resourceType":"video"}]}`)
	if response.Code != http.StatusOK || len(client.deletes) != 0 || !strings.Contains(response.Body.String(), "invalid delete request") {
		t.Fatalf("status = %d, deletes = %d, body = %s", response.Code, len(client.deletes), response.Body.String())
	}
}

func performDelete(t *testing.T, router http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, "/uploads", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
