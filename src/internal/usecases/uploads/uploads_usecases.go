package uploads

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/url"
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
	ErrInvalidClient = errors.New("upload asset client is not configured")

	allowedFolders = map[string]struct{}{
		"homestays":       {},
		"travel_location": {},
	}

	pathPattern = regexp.MustCompile(
		`^(homestays|travel_location)/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\.(jpg|png|webp|gif|mp4|webm|mov)$`,
	)
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

type UploadResult struct {
	URL string
}

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

// UseCase manages temporary uploaded assets.
//
// After a path has been persisted by a business workflow,
// that workflow must not call Delete for the persisted path.
type UseCase struct {
	client AssetClient
}

func NewUseCase(client AssetClient) *UseCase {
	return &UseCase{
		client: client,
	}
}

type validatedFile struct {
	File
	resourceType ResourceType
	mimeType     string
	extension    string
}

// Upload validates and uploads a batch of files.
//
// The maximum batch size is 10 files and the declared aggregate size
// must not exceed 50 MiB. Videos are limited to 10 MiB each.
//
// The HTTP handler must independently limit the actual request body size.
func (uc *UseCase) Upload(
	ctx context.Context,
	folder string,
	files []File,
) ([]UploadedFile, error) {
	if uc == nil || uc.client == nil {
		return nil, ErrInvalidClient
	}

	if _, ok := allowedFolders[folder]; !ok ||
		len(files) == 0 ||
		len(files) > maxFiles {
		return nil, ErrInvalidUpload
	}

	validated := make([]validatedFile, 0, len(files))
	var totalSize int64

	// Validate the entire batch before uploading any file.
	for _, file := range files {
		if file.Size <= 0 ||
			file.Size > maxRequestSize-totalSize ||
			file.Open == nil {
			return nil, ErrInvalidUpload
		}

		totalSize += file.Size

		kind, mimeType, extension, err := validateFile(file)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidUpload, err)
		}

		validated = append(validated, validatedFile{
			File:         file,
			resourceType: kind,
			mimeType:     mimeType,
			extension:    extension,
		})
	}

	uploaded := make([]UploadedFile, 0, len(validated))

	for _, file := range validated {
		if err := ctx.Err(); err != nil {
			uc.cleanup(ctx, uploaded)
			return nil, err
		}

		id, err := uuidV4()
		if err != nil {
			uc.cleanup(ctx, uploaded)
			return nil, err
		}

		reader, err := file.Open()
		if err != nil {
			uc.cleanup(ctx, uploaded)
			return nil, err
		}

		publicID := folder + "/" + id
		path := publicID + file.extension

		// Limit bytes exposed to the provider. The extra byte allows
		// detection of a file exceeding its declared size.
		limitedReader := &countingReader{
			reader: io.LimitReader(reader, file.Size+1),
		}

		result, uploadErr := uc.client.Upload(ctx, UploadRequest{
			Reader:       limitedReader,
			Folder:       folder,
			PublicID:     id,
			ResourceType: file.resourceType,
		})

		// Drain remaining bytes if the client didn't read the entire file
		_, _ = io.Copy(io.Discard, limitedReader)

		closeErr := reader.Close()

		// A provider may have created a resource even if Upload returned
		// an error. Attempt cleanup of the current resource as well.
		if uploadErr != nil {
			_ = uc.client.Delete(ctx, DeleteRequest{
				PublicID:     publicID,
				ResourceType: file.resourceType,
			})
			uc.cleanup(ctx, uploaded)
			return nil, uploadErr
		}

		if closeErr != nil {
			_ = uc.client.Delete(ctx, DeleteRequest{
				PublicID:     publicID,
				ResourceType: file.resourceType,
			})
			uc.cleanup(ctx, uploaded)
			return nil, closeErr
		}

		if limitedReader.n != file.Size {
			_ = uc.client.Delete(ctx, DeleteRequest{
				PublicID:     publicID,
				ResourceType: file.resourceType,
			})
			uc.cleanup(ctx, uploaded)
			return nil, fmt.Errorf(
				"%w: file size does not match declared size",
				ErrInvalidUpload,
			)
		}

		if !isValidHTTPSURL(result.URL) {
			_ = uc.client.Delete(ctx, DeleteRequest{
				PublicID:     publicID,
				ResourceType: file.resourceType,
			})
			uc.cleanup(ctx, uploaded)
			return nil, errors.New("upload provider returned an invalid URL")
		}

		uploaded = append(uploaded, UploadedFile{
			FileName: filepath.Base(file.FileName),
			MimeType: file.mimeType,
			Path:     path,
			URL:      result.URL,
		})
	}

	return uploaded, nil
}

type countingReader struct {
	reader io.Reader
	n      int64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.n += int64(n)
	return n, err
}

// Delete deletes assets independently and returns one result per item.
func (uc *UseCase) Delete(
	ctx context.Context,
	items []DeleteItem,
) []DeleteResult {
	results := make([]DeleteResult, len(items))

	if uc == nil || uc.client == nil {
		for i, item := range items {
			results[i] = DeleteResult{
				Path:    item.Path,
				Message: ErrInvalidClient.Error(),
			}
		}
		return results
	}

	for i, item := range items {
		if err := ctx.Err(); err != nil {
			results[i] = DeleteResult{
				Path:    item.Path,
				Message: err.Error(),
			}
			continue
		}

		publicID, resourceType, err := parseDeletePath(item)
		if err != nil {
			results[i] = DeleteResult{
				Path:    item.Path,
				Message: err.Error(),
			}
			continue
		}

		if err := uc.client.Delete(ctx, DeleteRequest{
			PublicID:     publicID,
			ResourceType: resourceType,
		}); err != nil {
			results[i] = DeleteResult{
				Path:    item.Path,
				Message: err.Error(),
			}
			continue
		}

		results[i] = DeleteResult{
			Path:    item.Path,
			Success: true,
		}
	}

	return results
}

// cleanup removes assets uploaded by the current, uncommitted batch.
// Cleanup uses a non-cancelled context so request cancellation does not
// immediately prevent deletion attempts.
func (uc *UseCase) cleanup(
	ctx context.Context,
	uploaded []UploadedFile,
) {
	if uc == nil || uc.client == nil {
		return
	}

	cleanupCtx := context.WithoutCancel(ctx)

	for _, file := range uploaded {
		matches := pathPattern.FindStringSubmatch(file.Path)
		if len(matches) != 4 {
			continue
		}

		_ = uc.client.Delete(cleanupCtx, DeleteRequest{
			PublicID: matches[1] + "/" + matches[2],
			ResourceType: resourceTypeForExtension(
				matches[3],
			),
		})
	}
}

func validateFile(
	file File,
) (ResourceType, string, string, error) {
	if file.Open == nil || file.Size <= 0 {
		return "", "", "", errors.New("invalid file")
	}

	reader, err := file.Open()
	if err != nil {
		return "", "", "", err
	}

	head := make([]byte, 12)
	n, readErr := io.ReadFull(reader, head)
	closeErr := reader.Close()

	if readErr != nil &&
		readErr != io.ErrUnexpectedEOF &&
		readErr != io.EOF {
		return "", "", "", readErr
	}

	if closeErr != nil {
		return "", "", "", closeErr
	}

	head = head[:n]

	switch {
	case len(head) >= 3 &&
		bytes.Equal(head[:3], []byte{0xff, 0xd8, 0xff}):
		return ResourceImage, "image/jpeg", ".jpg", nil

	case len(head) >= 8 &&
		bytes.Equal(
			head[:8],
			[]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a},
		):
		return ResourceImage, "image/png", ".png", nil

	case len(head) >= 6 &&
		(string(head[:6]) == "GIF87a" ||
			string(head[:6]) == "GIF89a"):
		return ResourceImage, "image/gif", ".gif", nil

	case len(head) >= 12 &&
		string(head[:4]) == "RIFF" &&
		string(head[8:12]) == "WEBP":
		return ResourceImage, "image/webp", ".webp", nil

	case len(head) >= 12 &&
		string(head[4:8]) == "ftyp":
		if file.Size > maxVideoSize {
			return "", "", "", errors.New("video exceeds 10 MiB")
		}

		brand := string(head[8:12])
		if brand == "qt  " {
			return ResourceVideo, "video/quicktime", ".mov", nil
		}

		return ResourceVideo, "video/mp4", ".mp4", nil

	case len(head) >= 4 &&
		bytes.Equal(head[:4], []byte{0x1a, 0x45, 0xdf, 0xa3}):
		if file.Size > maxVideoSize {
			return "", "", "", errors.New("video exceeds 10 MiB")
		}

		return ResourceVideo, "video/webm", ".webm", nil

	default:
		return "", "", "", errors.New(
			"unsupported or mismatched file content",
		)
	}
}

func parseDeletePath(
	item DeleteItem,
) (string, ResourceType, error) {
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
	switch extension {
	case "mp4", "webm", "mov":
		return ResourceVideo
	default:
		return ResourceImage
	}
}

func isValidHTTPSURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsed.Scheme, "https") &&
		parsed.Hostname() != "" &&
		parsed.User == nil
}

func uuidV4() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
