package controller

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/usecases/uploads"
)

type UploadController struct{ usecase *uploads.UseCase }

func NewUploadController(usecase *uploads.UseCase) *UploadController {
	return &UploadController{usecase: usecase}
}

func (uc *UploadController) UploadFiles(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50*1024*1024)
	if err := c.Request.ParseMultipartForm(32 * 1024 * 1024); err != nil {
		uc.writeUploadError(c, err)
		return
	}
	defer c.Request.MultipartForm.RemoveAll()
	files := c.Request.MultipartForm.File["files"]
	items, err := uc.usecase.Upload(c.Request.Context(), c.Request.FormValue("folder"), multipartFiles(files))
	if err != nil {
		uc.writeUploadError(c, err)
		return
	}
	response := make([]api.UploadItem, len(items))
	for i, item := range items {
		response[i] = api.UploadItem{FileName: item.FileName, MimeType: item.MimeType, Path: item.Path, Url: item.URL, RequiredHeaders: map[string]string{"Content-Type": item.MimeType}}
	}
	c.JSON(http.StatusOK, api.UploadResponse{Data: response})
}

func (uc *UploadController) DeleteUploads(c *gin.Context) {
	var body api.DeleteUploadsJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	items := make([]uploads.DeleteItem, len(body.Data))
	for i, item := range body.Data {
		items[i] = uploads.DeleteItem{Path: item.Path, ResourceType: uploads.ResourceType(item.ResourceType)}
	}
	results := uc.usecase.Delete(c.Request.Context(), items)
	response := make([]api.DeleteUploadResult, len(results))
	for i, result := range results {
		response[i] = api.DeleteUploadResult{Path: result.Path, Success: result.Success}
		if result.Message != "" {
			message := result.Message
			response[i].Message = &message
		}
	}
	c.JSON(http.StatusOK, api.DeleteUploadsResponse{Data: response})
}

func multipartFiles(headers []*multipart.FileHeader) []uploads.File {
	files := make([]uploads.File, len(headers))
	for i, header := range headers {
		fileHeader := header
		files[i] = uploads.File{FileName: fileHeader.Filename, Size: fileHeader.Size, Open: func() (io.ReadCloser, error) { return fileHeader.Open() }}
	}
	return files
}

func (uc *UploadController) writeUploadError(c *gin.Context, err error) {
	if errors.Is(err, uploads.ErrInvalidUpload) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusBadGateway, gin.H{"error": "upload provider failed"})
}
