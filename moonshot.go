// Package moonshot provides a complete SDK for the Moonshot AI API.
package moonshot

import (
	"github.com/rizome-dev/go-moonshot/pkg/chat"
	"github.com/rizome-dev/go-moonshot/pkg/client"
	"github.com/rizome-dev/go-moonshot/pkg/errors"
	"github.com/rizome-dev/go-moonshot/pkg/files"
	"github.com/rizome-dev/go-moonshot/pkg/models"
	"github.com/rizome-dev/go-moonshot/pkg/types"
	"github.com/rizome-dev/go-moonshot/pkg/utils"
)

// Re-export commonly used types for convenience
type (
	// Chat types
	ChatCompletionRequest  = types.ChatCompletionRequest
	ChatCompletionResponse = types.ChatCompletionResponse
	Message                = types.Message
	Choice                 = types.Choice
	Usage                  = types.Usage
	
	// Streaming types
	ChatCompletionStream       = types.ChatCompletionStream
	ChatCompletionStreamChoice = types.ChatCompletionStreamChoice
	ChatCompletionStreamDelta  = types.ChatCompletionStreamDelta
	
	// File types
	File           = types.File
	FileUploadReq  = types.FileUploadRequest
	FileListParams = types.FileListParams
	
	// Tool types
	Tool         = types.Tool
	ToolCall     = types.ToolCall
	ToolChoice   = types.ToolChoice
	Function     = types.Function
	FunctionCall = types.FunctionCall
	
	// Error types
	Error    = errors.Error
	APIError = errors.APIError
)

// Re-export model constants
var (
	ModelMoonshotV18K   = models.MoonshotV18K
	ModelMoonshotV132K  = models.MoonshotV132K
	ModelMoonshotV1128K = models.MoonshotV1128K
	ModelKimiK2         = models.KimiK2
	ModelKimiK2Base     = models.KimiK2Base
	ModelKimiK2Instruct = models.KimiK2Instruct
)

// Re-export utility functions
var (
	String  = utils.String
	Int     = utils.Int
	Int64   = utils.Int64
	Float64 = utils.Float64
	Float32 = utils.Float32
	Bool    = utils.Bool
	Time    = utils.Time
)

// Re-export error helper functions
var IsAPIError = errors.IsAPIError

// Re-export error code constants
const (
	ErrCodeInvalidRequest    = errors.ErrCodeInvalidRequest
	ErrCodeAuthentication    = errors.ErrCodeAuthentication
	ErrCodePermissionDenied  = errors.ErrCodePermissionDenied
	ErrCodeNotFound          = errors.ErrCodeNotFound
	ErrCodeRateLimitExceeded = errors.ErrCodeRateLimitExceeded
	ErrCodeServerError       = errors.ErrCodeServerError
	ErrCodeTimeout           = errors.ErrCodeTimeout
)

// SDK provides a convenient all-in-one client with all services
type SDK struct {
	Client *client.Client
	Chat   *chat.Service
	Files  *files.Service
}

// New creates a new Moonshot SDK instance with all services initialized.
// This is a convenience wrapper for users who want everything in one place.
//
// Usage:
//
//	sdk := moonshot.New()                           // Uses MOONSHOT_API_KEY env var
//	sdk := moonshot.New("sk-...")                  // Uses provided API key
//	sdk := moonshot.New(client.WithTimeout(30*time.Second))  // With options
func New(params ...interface{}) *SDK {
	c := client.New(params...)
	
	return &SDK{
		Client: c,
		Chat:   chat.NewService(c),
		Files:  files.NewService(c),
	}
}

// Version returns the SDK version
func Version() string {
	return "0.1.0"
}