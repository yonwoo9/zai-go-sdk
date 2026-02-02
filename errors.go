package zai

import (
	"fmt"
	"net/http"
)

// Error represents a ZAI API error
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	Type       string `json:"type,omitempty"`
	Code       string `json:"code,omitempty"`
}

func (e *Error) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("zai: %s (status: %d, type: %s, code: %s)", e.Message, e.StatusCode, e.Type, e.Code)
	}
	return fmt.Sprintf("zai: %s", e.Message)
}

// APIError represents a general API error
type APIError struct {
	Err *Error
}

func (e *APIError) Error() string { return e.Err.Error() }

// APIStatusError represents an API status error
type APIStatusError struct {
	Err *Error
}

func (e *APIStatusError) Error() string { return e.Err.Error() }

// APIRequestFailedError represents a 400 Bad Request error
type APIRequestFailedError struct {
	Err *Error
}

func (e *APIRequestFailedError) Error() string { return e.Err.Error() }

// APIAuthenticationError represents a 401 Unauthorized error
type APIAuthenticationError struct {
	Err *Error
}

func (e *APIAuthenticationError) Error() string { return e.Err.Error() }

// APIReachLimitError represents a 429 Rate Limit error
type APIReachLimitError struct {
	Err *Error
}

func (e *APIReachLimitError) Error() string { return e.Err.Error() }

// APIInternalError represents a 500 Internal Server error
type APIInternalError struct {
	Err *Error
}

func (e *APIInternalError) Error() string { return e.Err.Error() }

// APIServerFlowExceedError represents a 503 Service Unavailable error
type APIServerFlowExceedError struct {
	Err *Error
}

func (e *APIServerFlowExceedError) Error() string { return e.Err.Error() }

// APITimeoutError represents a timeout error
type APITimeoutError struct {
	Err *Error
}

func (e *APITimeoutError) Error() string { return e.Err.Error() }

// NewError creates a new Error based on HTTP status code
func NewError(statusCode int, message, errType, code string) error {
	baseErr := &Error{
		Message:    message,
		StatusCode: statusCode,
		Type:       errType,
		Code:       code,
	}

	switch statusCode {
	case http.StatusBadRequest:
		return &APIRequestFailedError{Err: baseErr}
	case http.StatusUnauthorized:
		return &APIAuthenticationError{Err: baseErr}
	case http.StatusTooManyRequests:
		return &APIReachLimitError{Err: baseErr}
	case http.StatusInternalServerError:
		return &APIInternalError{Err: baseErr}
	case http.StatusServiceUnavailable:
		return &APIServerFlowExceedError{Err: baseErr}
	default:
		return &APIStatusError{Err: baseErr}
	}
}
