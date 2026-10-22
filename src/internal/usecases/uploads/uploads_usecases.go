package uploads

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxFiles       = 10
	maxRequestSize = 50 * 1024 * 1024
	maxVideoSize   = 10 * 1024 * 1024
)

var (
	ErrInvalidUpload = errors.New("invalid upload")
	ErrInvalidDelete = errors.New("invalid delete request")
	allowedFolders   = map[string]struct{}{"homestays": {}, "travel": {}}
	pathPattern      = regexp.MustCompile(`^(homestays|travel)/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\.(jpg|png|webp|gif|mp4|webm|mov)$`)
)

type ResourceType string

const (
	ResourceImage ResourceType = "image"
	ResourceVideo ResourceType = "video"
)

type File struct {
	FileName string
	Size     int64
	Open     func() (io.ReadCloser, error)
}

type AssetClient interface {
	Upload(context.Context, UploadRequest) (UploadResult, error)
	Delete(context.Context, DeleteRequest) error
}

type UploadRequest struct {
	Reader       io.Reader
	Folder       string
	PublicID     string
	ResourceType ResourceType
}

type UploadResult struct{ URL string }

type DeleteRequest struct {
	PublicID     string
	ResourceType ResourceType
}

type UploadedFile struct {
	FileName string
	MimeType string
	Path     string
	URL      string
}

type DeleteItem struct {
	Path         string
	ResourceType ResourceType
}

type DeleteResult struct {
	Path    string
	Success bool
	Message string
}

// UseCase manages temporary uploaded assets. After a path has been persisted
// by a business workflow, that workflow must not call Delete for the path.
type UseCase struct{ client AssetClient }

func NewUseCase(client AssetClient) *UseCase { return &UseCase{client: client} }

func (uc *UseCase) Upload(ctx context.Context, folder string, files []File) ([]UploadedFile, error) {
	if _, ok := allowedFolders[folder]; !ok || len(files) == 0 || len(files) > maxFiles {
		return nil, ErrInvalidUpload
	}

	var totalSize int64
	validated := make([]validatedFile, len(files))
	for i, file := range files {
		totalSize += file.Size
		if totalSize > maxRequestSize || file.Size <= 0 {
			return nil, ErrInvalidUpload
		}
		kind, mimeType, extension, err := validateFile(file)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUpload, err)
		}
		validated[i] = validatedFile{File: file, resourceType: kind, mimeType: mimeType, extension: extension}
	}

	uploaded := make([]UploadedFile, 0, len(validated))
	for _, file := range validated {
		id, err := uuidV4()
		if err != nil {
			return nil, err
		}
		reader, err := file.Open()
		if err != nil {
			uc.cleanup(ctx, uploaded)
			return nil, err
		}
		result, uploadErr := uc.client.Upload(ctx, UploadRequest{Reader: reader, Folder: folder, PublicID: id, ResourceType: file.resourceType})
		closeErr := reader.Close()
		if uploadErr != nil || closeErr != nil || !strings.HasPrefix(result.URL, "https://") {
			uc.cleanup(ctx, uploaded)
			if uploadErr != nil {
				return nil, uploadErr
			}
			return nil, errors.New("upload provider returned an invalid URL")
		}
		uploaded = append(uploaded, UploadedFile{FileName: filepath.Base(file.FileName), MimeType: file.mimeType, Path: folder + "/" + id + file.extension, URL: result.URL})
	}
	return uploaded, nil
}

func (uc *UseCase) Delete(ctx context.Context, items []DeleteItem) []DeleteResult {
	results := make([]DeleteResult, len(items))
	for i, item := range items {
		publicID, resourceType, err := parseDeletePath(item)
		if err != nil {
			results[i] = DeleteResult{Path: item.Path, Message: err.Error()}
			continue
		}
		if err := uc.client.Delete(ctx, DeleteRequest{PublicID: publicID, ResourceType: resourceType}); err != nil {
			results[i] = DeleteResult{Path: item.Path, Message: err.Error()}
			continue
		}
		results[i] = DeleteResult{Path: item.Path, Success: true}
	}
	return results
}

func (uc *UseCase) cleanup(ctx context.Context, uploaded []UploadedFile) {
	for _, file := range uploaded {
		matches := pathPattern.FindStringSubmatch(file.Path)
		if len(matches) == 4 {
			_ = uc.client.Delete(ctx, DeleteRequest{PublicID: matches[1] + "/" + matches[2], ResourceType: resourceTypeForExtension(matches[3])})
		}
	}
}

type validatedFile struct {
	File
	resourceType ResourceType
	mimeType     string
	extension    string
}

func validateFile(file File) (ResourceType, string, string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", "", "", err
	}
	defer reader.Close()
	head := make([]byte, 32)
	n, err := io.ReadFull(reader, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", "", "", err
	}
	head = head[:n]
	switch {
	case len(head) >= 3 && bytes.Equal(head[:3], []byte{0xff, 0xd8, 0xff}):
		return ResourceImage, "image/jpeg", ".jpg", nil
	case len(head) >= 8 && bytes.Equal(head[:8], []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}):
		return ResourceImage, "image/png", ".png", nil
	case len(head) >= 6 && (string(head[:6]) == "GIF87a" || string(head[:6]) == "GIF89a"):
		return ResourceImage, "image/gif", ".gif", nil
	case len(head) >= 12 && string(head[:4]) == "RIFF" && string(head[8:12]) == "WEBP":
		return ResourceImage, "image/webp", ".webp", nil
	case len(head) >= 8 && string(head[4:8]) == "ftyp":
		if file.Size > maxVideoSize {
			return "", "", "", errors.New("video exceeds 10 MiB")
		}
		brand := string(head[8:12])
		if brand == "qt  " {
			return ResourceVideo, "video/quicktime", ".mov", nil
		}
		return ResourceVideo, "video/mp4", ".mp4", nil
	case len(head) >= 4 && bytes.Equal(head[:4], []byte{0x1a, 0x45, 0xdf, 0xa3}):
		if file.Size > maxVideoSize {
			return "", "", "", errors.New("video exceeds 10 MiB")
		}
		return ResourceVideo, "video/webm", ".webm", nil
	default:
		return "", "", "", errors.New("unsupported or mismatched file content")
	}
}

func parseDeletePath(item DeleteItem) (string, ResourceType, error) {
	matches := pathPattern.FindStringSubmatch(item.Path)
	if len(matches) != 4 {
		return "", "", ErrInvalidDelete
	}
	actualType := resourceTypeForExtension(matches[3])
	if item.ResourceType != actualType {
		return "", "", ErrInvalidDelete
	}
	return matches[1] + "/" + matches[2], actualType, nil
}

func resourceTypeForExtension(extension string) ResourceType {
	if extension == "mp4" || extension == "webm" || extension == "mov" {
		return ResourceVideo
	}
	return ResourceImage
}

func uuidV4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
