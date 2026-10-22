package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"platform-go/internal/usecases/uploads"
)

const requestTimeout = 30 * time.Second

type Client struct {
	cloudName string
	apiKey    string
	apiSecret string
	assetEnv  string
	assetErr  error
	http      *http.Client
}

func NewClient() *Client {
	assetEnv, assetErr := cloudinaryAssetEnvironment(os.Getenv("CLOUDINARY_ASSET_ENV"))
	return &Client{
		cloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		apiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		apiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		assetEnv:  assetEnv,
		assetErr:  assetErr,
		http:      &http.Client{Timeout: requestTimeout},
	}
}

func (c *Client) Upload(ctx context.Context, request uploads.UploadRequest) (uploads.UploadResult, error) {
	if err := c.configured(); err != nil {
		return uploads.UploadResult{}, err
	}
	timestamp := fmt.Sprint(time.Now().Unix())
	folder, err := c.scopedFolder(request.Folder)
	if err != nil {
		return uploads.UploadResult{}, err
	}
	params := map[string]string{"folder": folder, "public_id": request.PublicID, "timestamp": timestamp}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range params {
		_ = writer.WriteField(key, value)
	}
	_ = writer.WriteField("api_key", c.apiKey)
	_ = writer.WriteField("signature", c.sign(params))
	part, err := writer.CreateFormFile("file", request.PublicID)
	if err != nil {
		return uploads.UploadResult{}, err
	}
	if _, err := io.Copy(part, request.Reader); err != nil {
		return uploads.UploadResult{}, err
	}
	if err := writer.Close(); err != nil {
		return uploads.UploadResult{}, err
	}

	response, err := c.do(ctx, http.MethodPost, c.endpoint(request.ResourceType, "upload"), body, writer.FormDataContentType())
	if err != nil {
		return uploads.UploadResult{}, err
	}
	defer response.Body.Close()
	var payload struct {
		SecureURL string `json:"secure_url"`
		Error     struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return uploads.UploadResult{}, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return uploads.UploadResult{}, fmt.Errorf("cloudinary upload failed: %s", payload.Error.Message)
	}
	if !strings.HasPrefix(payload.SecureURL, "https://") {
		return uploads.UploadResult{}, errors.New("cloudinary returned no secure URL")
	}
	return uploads.UploadResult{URL: payload.SecureURL}, nil
}

func (c *Client) Delete(ctx context.Context, request uploads.DeleteRequest) error {
	if err := c.configured(); err != nil {
		return err
	}
	timestamp := fmt.Sprint(time.Now().Unix())
	publicID, err := c.scopedPublicID(request.PublicID)
	if err != nil {
		return err
	}
	params := map[string]string{"public_id": publicID, "timestamp": timestamp}
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	values.Set("api_key", c.apiKey)
	values.Set("signature", c.sign(params))
	response, err := c.do(ctx, http.MethodPost, c.endpoint(request.ResourceType, "destroy"), strings.NewReader(values.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var payload struct {
		Result string `json:"result"`
		Error  struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices || payload.Result != "ok" {
		if payload.Error.Message != "" {
			return errors.New(payload.Error.Message)
		}
		return fmt.Errorf("cloudinary delete failed: %s", payload.Result)
	}
	return nil
}

func (c *Client) configured() error {
	if c.cloudName == "" || c.apiKey == "" || c.apiSecret == "" {
		return errors.New("cloudinary is not configured")
	}
	_, err := c.environmentFolder()
	return err
}

func (c *Client) scopedFolder(folder string) (string, error) {
	environment, err := c.environmentFolder()
	if err != nil {
		return "", err
	}
	return path.Join(environment, folder), nil
}

func (c *Client) scopedPublicID(publicID string) (string, error) {
	environment, err := c.environmentFolder()
	if err != nil {
		return "", err
	}
	return path.Join(environment, publicID), nil
}

func (c *Client) environmentFolder() (string, error) {
	return c.assetEnv, c.assetErr
}

func cloudinaryAssetEnvironment(value string) (string, error) {
	switch value {
	case "dev", "stg":
		return value, nil
	default:
		return "", errors.New("CLOUDINARY_ASSET_ENV must be dev or stg")
	}
}
func (c *Client) endpoint(resourceType uploads.ResourceType, action string) string {
	return fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/%s", c.cloudName, resourceType, action)
}
func (c *Client) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+params[key])
	}
	hash := sha1.Sum([]byte(strings.Join(values, "&") + c.apiSecret))
	return hex.EncodeToString(hash[:])
}
func (c *Client) do(ctx context.Context, method, endpoint string, body io.Reader, contentType string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return c.http.Do(request)
}
