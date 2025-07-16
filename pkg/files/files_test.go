package files_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rizome-dev/go-moonshot/pkg/client"
	"github.com/rizome-dev/go-moonshot/pkg/errors"
	"github.com/rizome-dev/go-moonshot/pkg/files"
	"github.com/rizome-dev/go-moonshot/pkg/types"
)

func TestService_Upload(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		purpose  string
		handler  http.HandlerFunc
		want     *types.File
		wantErr  bool
	}{
		{
			name:     "successful upload",
			content:  "test file content",
			filename: "test.txt",
			purpose:  "assistants",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/files" {
					t.Errorf("expected /files, got %s", r.URL.Path)
				}
				
				// Parse multipart form
				err := r.ParseMultipartForm(32 << 20) // 32MB
				if err != nil {
					t.Errorf("failed to parse multipart form: %v", err)
				}
				
				// Verify purpose field
				purpose := r.FormValue("purpose")
				if purpose != "assistants" {
					t.Errorf("expected purpose 'assistants', got %s", purpose)
				}
				
				// Verify file
				file, header, err := r.FormFile("file")
				if err != nil {
					t.Errorf("failed to get file: %v", err)
				}
				defer file.Close()
				
				if header.Filename != "test.txt" {
					t.Errorf("expected filename 'test.txt', got %s", header.Filename)
				}
				
				// Read file content
				content, err := io.ReadAll(file)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
				}
				
				if string(content) != "test file content" {
					t.Errorf("expected content 'test file content', got %s", string(content))
				}
				
				// Send response
				resp := types.File{
					ID:        "file-123",
					Object:    "file",
					Bytes:     int64(len(content)),
					CreatedAt: 1234567890,
					Filename:  header.Filename,
					Purpose:   purpose,
					Status:    "processed",
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			want: &types.File{
				ID:        "file-123",
				Object:    "file",
				Bytes:     17,
				CreatedAt: 1234567890,
				Filename:  "test.txt",
				Purpose:   "assistants",
				Status:    "processed",
			},
			wantErr: false,
		},
		{
			name:     "API error",
			content:  "test",
			filename: "test.txt",
			purpose:  "invalid",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "invalid_request",
						Message: "Invalid purpose",
						Type:    "client_error",
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errResp)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			s := files.NewService(c)
			
			reader := strings.NewReader(tt.content)
			got, err := s.Upload(context.Background(), reader, tt.filename, tt.purpose)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Filename != tt.want.Filename {
					t.Errorf("Filename = %v, want %v", got.Filename, tt.want.Filename)
				}
				if got.Bytes != tt.want.Bytes {
					t.Errorf("Bytes = %v, want %v", got.Bytes, tt.want.Bytes)
				}
			}
		})
	}
}

func TestService_UploadFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	content := "test file content from disk"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Errorf("failed to parse multipart form: %v", err)
		}
		
		// Get the file
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("failed to get file: %v", err)
		}
		defer file.Close()
		
		// Verify filename matches the base name of the temp file
		expectedFilename := filepath.Base(tmpFile.Name())
		if header.Filename != expectedFilename {
			t.Errorf("expected filename %s, got %s", expectedFilename, header.Filename)
		}
		
		// Send response
		resp := types.File{
			ID:        "file-456",
			Object:    "file",
			Bytes:     int64(len(content)),
			CreatedAt: 1234567890,
			Filename:  header.Filename,
			Purpose:   r.FormValue("purpose"),
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	
	server := httptest.NewServer(handler)
	defer server.Close()
	
	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := files.NewService(c)
	
	got, err := s.UploadFile(context.Background(), tmpFile.Name(), "assistants")
	if err != nil {
		t.Errorf("UploadFile() error = %v", err)
	}
	
	if got != nil && got.ID != "file-456" {
		t.Errorf("ID = %v, want file-456", got.ID)
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name    string
		params  *types.FileListParams
		handler http.HandlerFunc
		want    *types.FileListResponse
		wantErr bool
	}{
		{
			name:   "list all files",
			params: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/files" {
					t.Errorf("expected /files, got %s", r.URL.Path)
				}
				
				resp := types.FileListResponse{
					Object: "list",
					Data: []types.File{
						{
							ID:       "file-1",
							Object:   "file",
							Filename: "file1.txt",
							Purpose:  "assistants",
						},
						{
							ID:       "file-2",
							Object:   "file",
							Filename: "file2.txt",
							Purpose:  "fine-tune",
						},
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			want: &types.FileListResponse{
				Object: "list",
				Data: []types.File{
					{ID: "file-1", Object: "file", Filename: "file1.txt", Purpose: "assistants"},
					{ID: "file-2", Object: "file", Filename: "file2.txt", Purpose: "fine-tune"},
				},
			},
			wantErr: false,
		},
		{
			name:   "list files with purpose filter",
			params: &types.FileListParams{Purpose: "assistants"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("purpose") != "assistants" {
					t.Errorf("expected purpose=assistants, got %s", r.URL.Query().Get("purpose"))
				}
				
				resp := types.FileListResponse{
					Object: "list",
					Data: []types.File{
						{ID: "file-1", Object: "file", Filename: "file1.txt", Purpose: "assistants"},
					},
				}
				
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			want: &types.FileListResponse{
				Object: "list",
				Data:   []types.File{{ID: "file-1", Object: "file", Filename: "file1.txt", Purpose: "assistants"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			s := files.NewService(c)
			
			got, err := s.List(context.Background(), tt.params)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if len(got.Data) != len(tt.want.Data) {
					t.Errorf("len(Data) = %v, want %v", len(got.Data), len(tt.want.Data))
				}
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name    string
		fileID  string
		handler http.HandlerFunc
		want    *types.File
		wantErr bool
	}{
		{
			name:   "get existing file",
			fileID: "file-123",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/files/file-123" {
					t.Errorf("expected /files/file-123, got %s", r.URL.Path)
				}
				
				resp := types.File{
					ID:        "file-123",
					Object:    "file",
					Bytes:     1024,
					CreatedAt: 1234567890,
					Filename:  "test.txt",
					Purpose:   "assistants",
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			want: &types.File{
				ID:        "file-123",
				Object:    "file",
				Bytes:     1024,
				CreatedAt: 1234567890,
				Filename:  "test.txt",
				Purpose:   "assistants",
			},
			wantErr: false,
		},
		{
			name:   "file not found",
			fileID: "file-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "not_found",
						Message: "File not found",
						Type:    "client_error",
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(errResp)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			s := files.NewService(c)
			
			got, err := s.Get(context.Background(), tt.fileID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Filename != tt.want.Filename {
					t.Errorf("Filename = %v, want %v", got.Filename, tt.want.Filename)
				}
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		fileID  string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:   "successful delete",
			fileID: "file-123",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Path != "/files/file-123" {
					t.Errorf("expected /files/file-123, got %s", r.URL.Path)
				}
				
				w.WriteHeader(http.StatusNoContent)
			},
			wantErr: false,
		},
		{
			name:   "delete with OK status",
			fileID: "file-456",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			},
			wantErr: false,
		},
		{
			name:   "file not found",
			fileID: "file-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "not_found",
						Message: "File not found",
						Type:    "client_error",
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(errResp)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			s := files.NewService(c)
			
			err := s.Delete(context.Background(), tt.fileID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetContent(t *testing.T) {
	tests := []struct {
		name    string
		fileID  string
		handler http.HandlerFunc
		want    []byte
		wantErr bool
	}{
		{
			name:   "get file content",
			fileID: "file-123",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/files/file-123/content" {
					t.Errorf("expected /files/file-123/content, got %s", r.URL.Path)
				}
				
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("This is the file content"))
			},
			want:    []byte("This is the file content"),
			wantErr: false,
		},
		{
			name:   "binary content",
			fileID: "file-bin",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte{0x00, 0x01, 0x02, 0x03})
			},
			want:    []byte{0x00, 0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name:   "file not found",
			fileID: "file-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "not_found",
						Message: "File not found",
						Type:    "client_error",
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(errResp)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			s := files.NewService(c)
			
			got, err := s.GetContent(context.Background(), tt.fileID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if !bytes.Equal(got, tt.want) {
					t.Errorf("content = %v, want %v", got, tt.want)
				}
			}
		})
	}
}