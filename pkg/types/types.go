package types

import (
	"encoding/json"
	"io"
)

// Message represents a message in a chat conversation
type Message struct {
	Role       string      `json:"role"`
	Content    interface{} `json:"content"`
	Name       *string     `json:"name,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID *string     `json:"tool_call_id,omitempty"`
}

// ContentPart represents a part of a message content (for multimodal)
type ContentPart struct {
	Type     string    `json:"type"`
	Text     *string   `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL in a message
type ImageURL struct {
	URL    string  `json:"url"`
	Detail *string `json:"detail,omitempty"`
}

// Tool represents a tool/function that can be called
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a function definition
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall represents a tool call in a message
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolChoice represents the tool choice configuration
type ToolChoice interface{}

// Usage represents token usage information
type Usage struct {
	PromptTokens       int     `json:"prompt_tokens"`
	CompletionTokens   int     `json:"completion_tokens"`
	TotalTokens        int     `json:"total_tokens"`
	PromptCacheHitRate float64 `json:"prompt_cache_hit_rate,omitempty"`
	PromptCacheMissRate float64 `json:"prompt_cache_miss_rate,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	FinishReason string   `json:"finish_reason"`
	LogProbs     *json.RawMessage `json:"logprobs,omitempty"`
}

// ChatCompletionRequest represents a request to the chat completion API
type ChatCompletionRequest struct {
	Model            string      `json:"model"`
	Messages         []Message   `json:"messages"`
	Temperature      *float64    `json:"temperature,omitempty"`
	TopP             *float64    `json:"top_p,omitempty"`
	N                *int        `json:"n,omitempty"`
	Stream           *bool       `json:"stream,omitempty"`
	Stop             interface{} `json:"stop,omitempty"`
	MaxTokens        *int        `json:"max_tokens,omitempty"`
	PresencePenalty  *float64    `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64    `json:"frequency_penalty,omitempty"`
	User             *string     `json:"user,omitempty"`
	
	// Tool use support
	Tools      []Tool     `json:"tools,omitempty"`
	ToolChoice ToolChoice `json:"tool_choice,omitempty"`
	
	// Moonshot-specific parameters
	ResponseFormat  interface{} `json:"response_format,omitempty"`
	Seed            *int64      `json:"seed,omitempty"`
	ParallelToolCalls *bool     `json:"parallel_tool_calls,omitempty"`
	
	// File references
	FileIDs []string `json:"file_ids,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completion API
type ChatCompletionResponse struct {
	ID                string    `json:"id"`
	Object            string    `json:"object"`
	Created           int64     `json:"created"`
	Model             string    `json:"model"`
	Choices           []Choice  `json:"choices"`
	Usage             Usage     `json:"usage"`
	SystemFingerprint string    `json:"system_fingerprint,omitempty"`
}

// ChatCompletionStream represents a streaming response chunk
type ChatCompletionStream struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	Choices           []ChatCompletionStreamChoice `json:"choices"`
	SystemFingerprint string                       `json:"system_fingerprint,omitempty"`
	Usage             *Usage                       `json:"usage,omitempty"`
}

// ChatCompletionStreamChoice represents a choice in a streaming response
type ChatCompletionStreamChoice struct {
	Index        int                       `json:"index"`
	Delta        ChatCompletionStreamDelta `json:"delta"`
	FinishReason *string                   `json:"finish_reason"`
	LogProbs     *json.RawMessage          `json:"logprobs,omitempty"`
}

// ChatCompletionStreamDelta represents the delta in a streaming response
type ChatCompletionStreamDelta struct {
	Role      *string    `json:"role,omitempty"`
	Content   *string    `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// File represents a file in the Moonshot API
type File struct {
	ID            string    `json:"id"`
	Object        string    `json:"object"`
	Bytes         int64     `json:"bytes"`
	CreatedAt     int64     `json:"created_at"`
	Filename      string    `json:"filename"`
	Purpose       string    `json:"purpose"`
	Status        string    `json:"status,omitempty"`
	StatusDetails *string   `json:"status_details,omitempty"`
}

// FileUploadRequest represents a request to upload a file
type FileUploadRequest struct {
	File    io.Reader `json:"-"`
	Purpose string    `json:"purpose"`
}

// FileListParams represents parameters for listing files
type FileListParams struct {
	Purpose string `json:"purpose,omitempty"`
}

// FileListResponse represents a response from listing files
type FileListResponse struct {
	Data   []File `json:"data"`
	Object string `json:"object"`
}

// TokenCountRequest represents a request to count tokens
type TokenCountRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// TokenCountResponse represents a response from token counting
type TokenCountResponse struct {
	TokenCount int `json:"token_count"`
}