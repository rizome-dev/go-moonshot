package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Error represents a general error
type Error struct {
	Message string
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Message
}

// APIError represents an error response from the Moonshot API
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("moonshot api error (status %d): %s - %s", e.StatusCode, e.Code, e.Message)
}

// ErrorResponse represents the structure of an error response from the API
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// HandleErrorResponse handles error responses from the API
func HandleErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIError{
			Code:       "read_error",
			Message:    fmt.Sprintf("failed to read error response: %v", err),
			Type:       "client_error",
			StatusCode: resp.StatusCode,
		}
	}
	
	// Try to parse as standard error response
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// If parsing fails, return a generic error with the body content
		return APIError{
			Code:       "parse_error",
			Message:    string(body),
			Type:       "client_error",
			StatusCode: resp.StatusCode,
		}
	}
	
	// Set the status code
	errResp.Error.StatusCode = resp.StatusCode
	
	return errResp.Error
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(APIError)
	if ok {
		return &apiErr, true
	}
	return nil, false
}

// Common error codes from Moonshot API
const (
	ErrCodeInvalidRequest     = "invalid_request"
	ErrCodeAuthentication     = "authentication_error"
	ErrCodePermissionDenied   = "permission_denied"
	ErrCodeNotFound          = "not_found"
	ErrCodeRateLimitExceeded = "rate_limit_exceeded"
	ErrCodeServerError       = "server_error"
	ErrCodeTimeout           = "timeout"
)