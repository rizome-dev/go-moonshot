package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/rizome-dev/go-moonshot/pkg/client"
	"github.com/rizome-dev/go-moonshot/pkg/errors"
	"github.com/rizome-dev/go-moonshot/pkg/types"
)

const (
	filesEndpoint = "/files"
)

// Service handles file-related operations
type Service struct {
	client *client.Client
}

// NewService creates a new files service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
	}
}

// Upload uploads a file to the Moonshot API
func (s *Service) Upload(ctx context.Context, file io.Reader, filename string, purpose string) (*types.File, error) {
	// Create a buffer to write our multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	
	// Create the file field
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", "application/octet-stream")
	
	fileWriter, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("creating file part: %w", err)
	}
	
	// Copy the file content
	if _, err := io.Copy(fileWriter, file); err != nil {
		return nil, fmt.Errorf("copying file content: %w", err)
	}
	
	// Add the purpose field
	if err := writer.WriteField("purpose", purpose); err != nil {
		return nil, fmt.Errorf("writing purpose field: %w", err)
	}
	
	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.BaseURL()+filesEndpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.client.APIKey())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Execute the request directly (not using client.Request which would override Content-Type)
	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	var fileResp types.File
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &fileResp, nil
}

// UploadFile is a convenience method that accepts a file path
func (s *Service) UploadFile(ctx context.Context, filePath string, purpose string) (*types.File, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()
	
	// Get the filename from the path
	filename := filepath.Base(filePath)
	
	return s.Upload(ctx, file, filename, purpose)
}

// List lists all files
func (s *Service) List(ctx context.Context, params *types.FileListParams) (*types.FileListResponse, error) {
	endpoint := filesEndpoint
	if params != nil && params.Purpose != "" {
		endpoint += "?purpose=" + params.Purpose
	}
	
	resp, err := s.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	var listResp types.FileListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &listResp, nil
}

// Get retrieves a file by ID
func (s *Service) Get(ctx context.Context, fileID string) (*types.File, error) {
	endpoint := fmt.Sprintf("%s/%s", filesEndpoint, fileID)
	
	resp, err := s.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	var fileResp types.File
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &fileResp, nil
}

// Delete deletes a file by ID
func (s *Service) Delete(ctx context.Context, fileID string) error {
	endpoint := fmt.Sprintf("%s/%s", filesEndpoint, fileID)
	
	resp, err := s.client.Request(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.HandleErrorResponse(resp)
	}
	
	return nil
}

// GetContent retrieves the content of a file
func (s *Service) GetContent(ctx context.Context, fileID string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s/content", filesEndpoint, fileID)
	
	resp, err := s.client.Request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading content: %w", err)
	}
	
	return content, nil
}