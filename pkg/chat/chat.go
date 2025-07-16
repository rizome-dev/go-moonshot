package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rizome-dev/go-moonshot/pkg/client"
	"github.com/rizome-dev/go-moonshot/pkg/errors"
	"github.com/rizome-dev/go-moonshot/pkg/types"
)

const (
	completionsEndpoint = "/chat/completions"
)

// Service handles chat-related operations
type Service struct {
	client *client.Client
}

// NewService creates a new chat service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
	}
}

// CreateCompletion creates a chat completion
func (s *Service) CreateCompletion(ctx context.Context, req types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {
	// Ensure streaming is disabled for non-streaming request
	req.Stream = &[]bool{false}[0]
	
	// Adjust temperature for Moonshot API (maps by real_temperature = request_temperature * 0.6)
	if req.Temperature != nil {
		adjustedTemp := *req.Temperature * 0.6
		req.Temperature = &adjustedTemp
	}
	
	resp, err := s.client.Request(ctx, http.MethodPost, completionsEndpoint, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	var completionResp types.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &completionResp, nil
}

// StreamReader represents a reader for streaming responses
type StreamReader struct {
	reader   *bufio.Reader
	response *http.Response
}

// Read reads the next streaming chunk
func (sr *StreamReader) Read() (*types.ChatCompletionStream, error) {
	line, err := sr.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("reading stream: %w", err)
	}
	
	line = strings.TrimSpace(line)
	
	// Skip empty lines
	if line == "" {
		return sr.Read()
	}
	
	// Check for data prefix
	if !strings.HasPrefix(line, "data: ") {
		return sr.Read()
	}
	
	// Remove "data: " prefix
	data := strings.TrimPrefix(line, "data: ")
	
	// Check for end of stream
	if data == "[DONE]" {
		return nil, io.EOF
	}
	
	// Parse the JSON
	var chunk types.ChatCompletionStream
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return nil, fmt.Errorf("parsing stream chunk: %w", err)
	}
	
	return &chunk, nil
}

// Close closes the stream reader
func (sr *StreamReader) Close() error {
	if sr.response != nil {
		return sr.response.Body.Close()
	}
	return nil
}

// CreateCompletionStream creates a streaming chat completion
func (s *Service) CreateCompletionStream(ctx context.Context, req types.ChatCompletionRequest) (*StreamReader, error) {
	// Ensure streaming is enabled
	req.Stream = &[]bool{true}[0]
	
	// Adjust temperature for Moonshot API (maps by real_temperature = request_temperature * 0.6)
	if req.Temperature != nil {
		adjustedTemp := *req.Temperature * 0.6
		req.Temperature = &adjustedTemp
	}
	
	resp, err := s.client.StreamRequest(ctx, http.MethodPost, completionsEndpoint, req)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, errors.HandleErrorResponse(resp)
	}
	
	return &StreamReader{
		reader:   bufio.NewReader(resp.Body),
		response: resp,
	}, nil
}

// CreateCompletionWithCallback creates a streaming chat completion with a callback for each chunk
func (s *Service) CreateCompletionWithCallback(ctx context.Context, req types.ChatCompletionRequest, callback func(*types.ChatCompletionStream) error) error {
	stream, err := s.CreateCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()
	
	for {
		chunk, err := stream.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		
		if err := callback(chunk); err != nil {
			return err
		}
	}
}

// CountTokens counts the number of tokens in a message sequence
func (s *Service) CountTokens(ctx context.Context, req types.TokenCountRequest) (*types.TokenCountResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/tokenizers/estimate_token_count", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.HandleErrorResponse(resp)
	}
	
	var tokenResp types.TokenCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &tokenResp, nil
}