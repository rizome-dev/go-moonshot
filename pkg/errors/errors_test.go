package errors_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/rizome-dev/go-moonshot/pkg/errors"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  errors.Error
		want string
	}{
		{
			name: "simple error",
			err:  errors.Error{Message: "something went wrong"},
			want: "something went wrong",
		},
		{
			name: "empty message",
			err:  errors.Error{Message: ""},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  errors.APIError
		want string
	}{
		{
			name: "complete error",
			err: errors.APIError{
				Code:       "invalid_request",
				Message:    "The request is invalid",
				Type:       "client_error",
				StatusCode: 400,
			},
			want: "moonshot api error (status 400): invalid_request - The request is invalid",
		},
		{
			name: "error without status code",
			err: errors.APIError{
				Code:    "server_error",
				Message: "Internal server error",
				Type:    "server_error",
			},
			want: "moonshot api error (status 0): server_error - Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    errors.APIError
	}{
		{
			name:       "standard error response",
			statusCode: 400,
			body:       `{"error": {"code": "invalid_request", "message": "Invalid API key", "type": "authentication_error"}}`,
			wantErr: errors.APIError{
				Code:       "invalid_request",
				Message:    "Invalid API key",
				Type:       "authentication_error",
				StatusCode: 400,
			},
		},
		{
			name:       "malformed JSON response",
			statusCode: 500,
			body:       `{invalid json`,
			wantErr: errors.APIError{
				Code:       "parse_error",
				Message:    "{invalid json",
				Type:       "client_error",
				StatusCode: 500,
			},
		},
		{
			name:       "empty response",
			statusCode: 503,
			body:       "",
			wantErr: errors.APIError{
				Code:       "parse_error",
				Message:    "",
				Type:       "client_error",
				StatusCode: 503,
			},
		},
		{
			name:       "plain text error",
			statusCode: 404,
			body:       "Not Found",
			wantErr: errors.APIError{
				Code:       "parse_error",
				Message:    "Not Found",
				Type:       "client_error",
				StatusCode: 404,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
			}
			
			err := errors.HandleErrorResponse(resp)
			
			apiErr, ok := err.(errors.APIError)
			if !ok {
				t.Fatalf("HandleErrorResponse() returned non-APIError type: %T", err)
			}
			
			if apiErr.Code != tt.wantErr.Code {
				t.Errorf("APIError.Code = %v, want %v", apiErr.Code, tt.wantErr.Code)
			}
			if apiErr.Message != tt.wantErr.Message {
				t.Errorf("APIError.Message = %v, want %v", apiErr.Message, tt.wantErr.Message)
			}
			if apiErr.Type != tt.wantErr.Type {
				t.Errorf("APIError.Type = %v, want %v", apiErr.Type, tt.wantErr.Type)
			}
			if apiErr.StatusCode != tt.wantErr.StatusCode {
				t.Errorf("APIError.StatusCode = %v, want %v", apiErr.StatusCode, tt.wantErr.StatusCode)
			}
		})
	}
}

func TestIsAPIError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantOk  bool
		wantErr *errors.APIError
	}{
		{
			name: "is API error",
			err: errors.APIError{
				Code:    "test_error",
				Message: "Test error",
			},
			wantOk: true,
			wantErr: &errors.APIError{
				Code:    "test_error",
				Message: "Test error",
			},
		},
		{
			name:    "is not API error",
			err:     errors.Error{Message: "generic error"},
			wantOk:  false,
			wantErr: nil,
		},
		{
			name:    "nil error",
			err:     nil,
			wantOk:  false,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotOk := errors.IsAPIError(tt.err)
			
			if gotOk != tt.wantOk {
				t.Errorf("IsAPIError() ok = %v, want %v", gotOk, tt.wantOk)
			}
			
			if gotOk && tt.wantOk {
				if gotErr.Code != tt.wantErr.Code || gotErr.Message != tt.wantErr.Message {
					t.Errorf("IsAPIError() error = %+v, want %+v", gotErr, tt.wantErr)
				}
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Just verify that the constants are defined
	constants := []string{
		errors.ErrCodeInvalidRequest,
		errors.ErrCodeAuthentication,
		errors.ErrCodePermissionDenied,
		errors.ErrCodeNotFound,
		errors.ErrCodeRateLimitExceeded,
		errors.ErrCodeServerError,
		errors.ErrCodeTimeout,
	}
	
	expectedValues := []string{
		"invalid_request",
		"authentication_error",
		"permission_denied",
		"not_found",
		"rate_limit_exceeded",
		"server_error",
		"timeout",
	}
	
	for i, constant := range constants {
		if constant != expectedValues[i] {
			t.Errorf("Error constant %d = %v, want %v", i, constant, expectedValues[i])
		}
	}
}