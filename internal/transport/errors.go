package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode represents API error codes.
type ErrorCode string

// Known error codes.
const (
	ErrCodeMissingParam      ErrorCode = "MISSING_REQUIRED_PARAMETER"
	ErrCodePromptTooLong     ErrorCode = "PROMPT_TOO_LONG"
	ErrCodeContentViolation  ErrorCode = "CONTENT_POLICY_VIOLATION"
	ErrCodeIndexOutOfBounds  ErrorCode = "INDEX_OUT_OF_BOUNDS"
	ErrCodeInvalidAPIKey     ErrorCode = "INVALID_API_KEY"
	ErrCodeInsufficientFunds ErrorCode = "INSUFFICIENT_CREDITS"
	ErrCodeRateLimit         ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternal          ErrorCode = "INTERNAL_ERROR"
)

// APIError represents an API error.
type APIError struct {
	Code       ErrorCode      `json:"error_code"`
	Message    string         `json:"message"`
	Params     map[string]any `json:"params,omitempty"`
	StatusCode int            `json:"-"`
	RequestID  string         `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("reve: %s (code=%s, status=%d, request_id=%s)",
			e.Message, e.Code, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("reve: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
}

// Retryable returns true if the error can be retried.
func (e *APIError) Retryable() bool {
	return isRetryableStatus(e.StatusCode)
}

// IsRateLimit returns true if rate limited.
func (e *APIError) IsRateLimit() bool {
	return e.Code == ErrCodeRateLimit || e.StatusCode == http.StatusTooManyRequests
}

// IsInsufficientFunds returns true if insufficient credits.
func (e *APIError) IsInsufficientFunds() bool {
	return e.Code == ErrCodeInsufficientFunds || e.StatusCode == http.StatusPaymentRequired
}

// IsContentViolation returns true if content policy violated.
func (e *APIError) IsContentViolation() bool {
	return e.Code == ErrCodeContentViolation
}

// IsAuthError returns true if authentication failed.
func (e *APIError) IsAuthError() bool {
	return e.Code == ErrCodeInvalidAPIKey || e.StatusCode == http.StatusUnauthorized
}

// RequestError represents a request-level error.
type RequestError struct {
	Op  string
	Err error
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	return fmt.Sprintf("reve: %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// ParseError parses an error response.
func ParseError(resp *http.Response, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Reve-Request-Id"),
	}

	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Message = string(body)
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
	}

	if apiErr.Code == "" {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			apiErr.Code = ErrCodeInvalidAPIKey
		case http.StatusPaymentRequired:
			apiErr.Code = ErrCodeInsufficientFunds
		case http.StatusTooManyRequests:
			apiErr.Code = ErrCodeRateLimit
		case http.StatusInternalServerError:
			apiErr.Code = ErrCodeInternal
		}
	}

	return apiErr
}
