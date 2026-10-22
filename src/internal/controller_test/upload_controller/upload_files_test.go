package upload_controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/usecases/uploads"
)

type fakeAssetClient struct {
	uploads []uploads.UploadRequest
	deletes []uploads.DeleteRequest
	err     error
}

func (c *fakeAssetClient) Upload(_ context.Context, request uploads.UploadRequest) (uploads.UploadResult, error) {
	c.uploads = append(c.uploads, request)
	if c.err != nil {
		return uploads.UploadResult{}, c.err
	}
	return uploads.UploadResult{URL: "https://res.cloudinary.com/example/asset"}, nil
}

func (c *fakeAssetClient) Delete(_ context.Context, request uploads.DeleteRequest) error {
	c.deletes = append(c.deletes, request)
	return c.err
}

func newUploadRouter(t *testing.T, client uploads.AssetClient) *gin.Engine {
	t.Helper()
	usecase := uploads.NewUseCase(client)
	uploadController := controller.NewUploadController(usecase)
	router := shared.NewSetupRouter(t)
	router.POST("/uploads", uploadController.UploadFiles)
	router.DELETE("/uploads", uploadController.DeleteUploads)
	return router
}

func TestUploadFiles(t *testing.T) {
	cases := []struct {
		name          string
		fileName      string
		content       []byte
		wantMIME      string
		wantExtension string
		wantType      uploads.ResourceType
	}{
		{"PNG", "spoofed.exe", pngContent(), "image/png", ".png", uploads.ResourceImage},
		{"JPEG", "room.png", []byte{0xff, 0xd8, 0xff, 0xdb}, "image/jpeg", ".jpg", uploads.ResourceImage},
		{"GIF", "animation.jpg", []byte("GIF89a"), "image/gif", ".gif", uploads.ResourceImage},
		{"MP4", "movie.gif", mp4Content(), "video/mp4", ".mp4", uploads.ResourceVideo},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeAssetClient{}
			response := performUpload(t, newUploadRouter(t, client), "homestays", test.fileName, test.content)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}

			var body api.UploadResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if len(body.Data) != 1 || body.Data[0].MimeType != test.wantMIME || !strings.HasSuffix(body.Data[0].Path, test.wantExtension) {
				t.Fatalf("response = %#v", body)
			}
			if body.Data[0].RequiredHeaders["Content-Type"] != test.wantMIME || len(client.uploads) != 1 || client.uploads[0].ResourceType != test.wantType {
				t.Fatalf("upload request = %#v, response = %#v", client.uploads, body.Data[0])
			}
		})
	}
}

func TestUploadFilesValidation(t *testing.T) {
	t.Run("accepts a 10 MiB video", func(t *testing.T) {
		content := make([]byte, 10*1024*1024)
		copy(content, mp4Content())
		response := performUpload(t, newUploadRouter(t, &fakeAssetClient{}), "travel", "video.mp4", content)
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
		}
	})

	cases := []struct {
		name, folder, fileName string
		content                []byte
	}{
		{"rejects video over 10 MiB", "travel", "video.mp4", oversizedMP4Content()},
		{"rejects untrusted folder", "../outside", "image.png", pngContent()},
		{"rejects spoofed MIME", "homestays", "image.png", []byte("not an image")},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeAssetClient{}
			response := performUpload(t, newUploadRouter(t, client), test.folder, test.fileName, test.content)
			if response.Code != http.StatusBadRequest || len(client.uploads) != 0 {
				t.Fatalf("status = %d, uploads = %d, body = %s", response.Code, len(client.uploads), response.Body.String())
			}
		})
	}
}

func TestUploadFilesReturnsProviderTimeout(t *testing.T) {
	response := performUpload(t, newUploadRouter(t, &fakeAssetClient{err: context.DeadlineExceeded}), "homestays", "image.png", pngContent())
	if response.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func performUpload(t *testing.T, router *gin.Engine, folder, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("folder", folder); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/uploads", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func pngContent() []byte { return []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a} }
func mp4Content() []byte { return []byte{0, 0, 0, 24, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm'} }

func oversizedMP4Content() []byte {
	content := make([]byte, 10*1024*1024+1)
	copy(content, mp4Content())
	return content
}
