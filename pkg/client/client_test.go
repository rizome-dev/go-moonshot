package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rizome-dev/go-moonshot/pkg/client"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		params []interface{}
		envKey string
		want   struct {
			hasAPIKey bool
			baseURL   string
		}
	}{
		{
			name:   "with api key string",
			params: []interface{}{"test-api-key"},
			want: struct {
				hasAPIKey bool
				baseURL   string
			}{
				hasAPIKey: true,
				baseURL:   "https://api.moonshot.ai/v1",
			},
		},
		{
			name:   "with env var",
			envKey: "test-env-key",
			params: []interface{}{},
			want: struct {
				hasAPIKey bool
				baseURL   string
			}{
				hasAPIKey: true,
				baseURL:   "https://api.moonshot.ai/v1",
			},
		},
		{
			name: "with custom base URL",
			params: []interface{}{
				"test-api-key",
				client.WithBaseURL("https://custom.api.com"),
			},
			want: struct {
				hasAPIKey bool
				baseURL   string
			}{
				hasAPIKey: true,
				baseURL:   "https://custom.api.com",
			},
		},
		{
			name: "with custom timeout",
			params: []interface{}{
				"test-api-key",
				client.WithTimeout(60 * time.Second),
			},
			want: struct {
				hasAPIKey bool
				baseURL   string
			}{
				hasAPIKey: true,
				baseURL:   "https://api.moonshot.ai/v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				os.Setenv(client.EnvAPIKey, tt.envKey)
				defer os.Unsetenv(client.EnvAPIKey)
			}
			
			c := client.New(tt.params...)
			
			if c == nil {
				t.Fatal("New() returned nil")
			}
			
			if tt.want.hasAPIKey && c.APIKey() == "" {
				t.Error("expected API key to be set")
			}
			
			if c.BaseURL() != tt.want.baseURL {
				t.Errorf("BaseURL() = %v, want %v", c.BaseURL(), tt.want.baseURL)
			}
		})
	}
}

func TestClient_Request(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:   "successful GET request",
			method: http.MethodGet,
			path:   "/test",
			body:   nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method %v, got %v", http.MethodGet, r.Method)
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("missing or incorrect Authorization header")
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
			},
			wantErr: false,
		},
		{
			name:   "successful POST request with body",
			method: http.MethodPost,
			path:   "/test",
			body: map[string]string{
				"message": "hello",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected method %v, got %v", http.MethodPost, r.Method)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing or incorrect Content-Type header")
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
			},
			wantErr: false,
		},
		{
			name:   "network error",
			method: http.MethodGet,
			path:   "/test",
			body:   nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Close connection to simulate network error
				hj, _ := w.(http.Hijacker)
				conn, _, _ := hj.Hijack()
				conn.Close()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			
			c := client.New("test-key", client.WithBaseURL(server.URL))
			
			resp, err := c.Request(context.Background(), tt.method, tt.path, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if resp != nil {
				defer resp.Body.Close()
			}
		})
	}
}

func TestOptions(t *testing.T) {
	t.Run("WithHTTPClient", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 100 * time.Second,
		}
		c := client.New("test-key", client.WithHTTPClient(customClient))
		if c == nil {
			t.Fatal("New() returned nil")
		}
		// Note: We can't directly test the HTTP client as it's not exported
	})
	
	t.Run("WithUserAgent", func(t *testing.T) {
		c := client.New("test-key", client.WithUserAgent("custom-agent/1.0"))
		if c == nil {
			t.Fatal("New() returned nil")
		}
		// Note: We can't directly test the user agent as it's not exported
	})
}