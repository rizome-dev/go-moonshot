package chat_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rizome-dev/go-moonshot/pkg/chat"
	"github.com/rizome-dev/go-moonshot/pkg/client"
	"github.com/rizome-dev/go-moonshot/pkg/errors"
	"github.com/rizome-dev/go-moonshot/pkg/models"
	"github.com/rizome-dev/go-moonshot/pkg/types"
	"github.com/rizome-dev/go-moonshot/pkg/utils"
)

func TestService_CreateCompletion(t *testing.T) {
	tests := []struct {
		name    string
		req     types.ChatCompletionRequest
		handler http.HandlerFunc
		want    *types.ChatCompletionResponse
		wantErr bool
	}{
		{
			name: "successful completion",
			req: types.ChatCompletionRequest{
				Model: models.MoonshotV18K.String(),
				Messages: []types.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/chat/completions" {
					t.Errorf("expected /chat/completions, got %s", r.URL.Path)
				}
				
				// Decode request body
				var req types.ChatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}
				
				// Verify stream is false
				if req.Stream == nil || *req.Stream != false {
					t.Error("expected stream to be false")
				}
				
				// Send response
				resp := types.ChatCompletionResponse{
					ID:      "chatcmpl-123",
					Object:  "chat.completion",
					Created: 1234567890,
					Model:   models.MoonshotV18K.String(),
					Choices: []types.Choice{
						{
							Index: 0,
							Message: types.Message{
								Role:    "assistant",
								Content: "Hi there!",
							},
							FinishReason: "stop",
						},
					},
					Usage: types.Usage{
						PromptTokens:     10,
						CompletionTokens: 5,
						TotalTokens:      15,
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			want: &types.ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: 1234567890,
				Model:   models.MoonshotV18K.String(),
				Choices: []types.Choice{
					{
						Index: 0,
						Message: types.Message{
							Role:    "assistant",
							Content: "Hi there!",
						},
						FinishReason: "stop",
					},
				},
				Usage: types.Usage{
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
				},
			},
			wantErr: false,
		},
		{
			name: "with temperature adjustment",
			req: types.ChatCompletionRequest{
				Model: models.MoonshotV18K.String(),
				Messages: []types.Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: utils.Float64(1.0),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Decode request body
				var req types.ChatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}
				
				// Verify temperature was adjusted
				if req.Temperature == nil || *req.Temperature != 0.6 {
					t.Errorf("expected temperature to be 0.6, got %v", req.Temperature)
				}
				
				// Send response
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(types.ChatCompletionResponse{
					ID:      "chatcmpl-123",
					Object:  "chat.completion",
					Created: 1234567890,
					Model:   models.MoonshotV18K.String(),
					Choices: []types.Choice{{Index: 0, Message: types.Message{Role: "assistant", Content: "Hi"}, FinishReason: "stop"}},
				})
			},
			want: &types.ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: 1234567890,
				Model:   models.MoonshotV18K.String(),
				Choices: []types.Choice{{Index: 0, Message: types.Message{Role: "assistant", Content: "Hi"}, FinishReason: "stop"}},
			},
			wantErr: false,
		},
		{
			name: "API error",
			req: types.ChatCompletionRequest{
				Model: models.MoonshotV18K.String(),
				Messages: []types.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "invalid_request",
						Message: "Invalid model",
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
			s := chat.NewService(c)
			
			got, err := s.CreateCompletion(context.Background(), tt.req)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				// Compare key fields
				if got.ID != tt.want.ID {
					t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
				}
				if len(got.Choices) != len(tt.want.Choices) {
					t.Errorf("len(Choices) = %v, want %v", len(got.Choices), len(tt.want.Choices))
				}
			}
		})
	}
}

func TestService_CreateCompletionStream(t *testing.T) {
	tests := []struct {
		name    string
		req     types.ChatCompletionRequest
		handler http.HandlerFunc
		chunks  []string
		wantErr bool
	}{
		{
			name: "successful streaming",
			req: types.ChatCompletionRequest{
				Model: models.MoonshotV18K.String(),
				Messages: []types.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				var req types.ChatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}
				
				// Verify stream is true
				if req.Stream == nil || *req.Stream != true {
					t.Error("expected stream to be true")
				}
				
				// Send SSE response
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)
				
				// Simulate streaming chunks
				chunks := []string{
					`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}`,
					`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
					`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{"content":" there!"},"finish_reason":null}]}`,
					`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
				}
				
				for _, chunk := range chunks {
					fmt.Fprintf(w, "data: %s\n\n", chunk)
					w.(http.Flusher).Flush()
				}
				
				fmt.Fprintf(w, "data: [DONE]\n\n")
				w.(http.Flusher).Flush()
			},
			chunks: []string{
				"Hello",
				" there!",
			},
			wantErr: false,
		},
		{
			name: "API error on stream",
			req: types.ChatCompletionRequest{
				Model: models.MoonshotV18K.String(),
				Messages: []types.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				errResp := errors.ErrorResponse{
					Error: errors.APIError{
						Code:    "rate_limit_exceeded",
						Message: "Rate limit exceeded",
						Type:    "client_error",
					},
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
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
			s := chat.NewService(c)
			
			stream, err := s.CreateCompletionStream(context.Background(), tt.req)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCompletionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && stream != nil {
				defer stream.Close()
				
				// Read chunks
				var contents []string
				for {
					chunk, err := stream.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Errorf("Read() error = %v", err)
						break
					}
					
					if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
						contents = append(contents, *chunk.Choices[0].Delta.Content)
					}
				}
				
				// Verify we got expected content
				if len(contents) != len(tt.chunks) {
					t.Errorf("got %d chunks, want %d", len(contents), len(tt.chunks))
				}
				for i, content := range contents {
					if i < len(tt.chunks) && content != tt.chunks[i] {
						t.Errorf("chunk[%d] = %v, want %v", i, content, tt.chunks[i])
					}
				}
			}
		})
	}
}

func TestService_CreateCompletionWithCallback(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		
		chunks := []string{
			`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{"content":"Test"},"finish_reason":null}]}`,
			`{"id":"123","object":"chat.completion.chunk","created":1234567890,"model":"moonshot-v1-8k","choices":[{"index":0,"delta":{"content":" message"},"finish_reason":null}]}`,
		}
		
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
			w.(http.Flusher).Flush()
		}
		
		fmt.Fprintf(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	})
	
	server := httptest.NewServer(handler)
	defer server.Close()
	
	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := chat.NewService(c)
	
	req := types.ChatCompletionRequest{
		Model: models.MoonshotV18K.String(),
		Messages: []types.Message{
			{Role: "user", Content: "Hello"},
		},
	}
	
	var receivedChunks []string
	callback := func(chunk *types.ChatCompletionStream) error {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
			receivedChunks = append(receivedChunks, *chunk.Choices[0].Delta.Content)
		}
		return nil
	}
	
	err := s.CreateCompletionWithCallback(context.Background(), req, callback)
	if err != nil {
		t.Errorf("CreateCompletionWithCallback() error = %v", err)
	}
	
	expectedChunks := []string{"Test", " message"}
	if len(receivedChunks) != len(expectedChunks) {
		t.Errorf("got %d chunks, want %d", len(receivedChunks), len(expectedChunks))
	}
}

func TestStreamReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantEOF bool
	}{
		{
			name: "parse valid SSE data",
			input: `data: {"id":"123","choices":[{"delta":{"content":"Hello"}}]}

data: {"id":"123","choices":[{"delta":{"content":" world"}}]}

data: [DONE]

`,
			want:    []string{"Hello", " world"},
			wantEOF: true,
		},
		{
			name: "skip empty lines",
			input: `

data: {"id":"123","choices":[{"delta":{"content":"Test"}}]}


data: [DONE]
`,
			want:    []string{"Test"},
			wantEOF: true,
		},
		{
			name: "skip non-data lines",
			input: `event: message
data: {"id":"123","choices":[{"delta":{"content":"Test"}}]}
comment: this is a comment

data: [DONE]
`,
			want:    []string{"Test"},
			wantEOF: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			
			// Simulate reading from SSE stream
			var contents []string
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Errorf("ReadString() error = %v", err)
					break
				}
				
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				if !strings.HasPrefix(line, "data: ") {
					continue
				}
				
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					break
				}
				
				var chunk types.ChatCompletionStream
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					t.Errorf("Unmarshal() error = %v", err)
					continue
				}
				
				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
					contents = append(contents, *chunk.Choices[0].Delta.Content)
				}
			}
			
			if len(contents) != len(tt.want) {
				t.Errorf("got %d contents, want %d", len(contents), len(tt.want))
			}
			
			for i, content := range contents {
				if i < len(tt.want) && content != tt.want[i] {
					t.Errorf("content[%d] = %v, want %v", i, content, tt.want[i])
				}
			}
		})
	}
}