package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL = "https://api.moonshot.ai/v1"
	defaultTimeout = 30 * time.Second
	EnvAPIKey      = "MOONSHOT_API_KEY"
	Version        = "0.1.0"
)

// Client is the main client for interacting with the Moonshot API
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
}

// Option is a function that configures a Client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithUserAgent sets a custom user agent
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// New creates a new Moonshot client.
// Usage:
//
//	client.New()                    // Uses MOONSHOT_API_KEY env var
//	client.New("sk-...")            // Uses provided API key
//	client.New(client.WithTimeout(30*time.Second))  // With options only
//	client.New("sk-...", client.WithTimeout(30*time.Second))  // With API key and options
func New(params ...interface{}) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL:   defaultBaseURL,
		userAgent: fmt.Sprintf("go-moonshot/%s", Version),
	}
	
	// Process parameters - can be string (API key) or Option
	for _, param := range params {
		switch v := param.(type) {
		case string:
			c.apiKey = v
		case Option:
			v(c)
		case func(*Client):
			v(c)
		}
	}
	
	// If no API key provided, try environment variable
	if c.apiKey == "" {
		c.apiKey = os.Getenv(EnvAPIKey)
	}
	
	return c
}

// Request performs an HTTP request to the Moonshot API
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path
	
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	
	return resp, nil
}

// StreamRequest performs a streaming HTTP request to the Moonshot API
func (c *Client) StreamRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// For streaming requests, we need to ensure the body includes stream: true
	// This is handled by the caller, but we use the same request method
	return c.Request(ctx, method, path, body)
}

// BaseURL returns the base URL of the client
func (c *Client) BaseURL() string {
	return c.baseURL
}

// APIKey returns the API key (useful for services that need it)
func (c *Client) APIKey() string {
	return c.apiKey
}