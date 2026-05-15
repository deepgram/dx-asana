package asana

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// APIError is a structured error returned by Asana's API. It implements the
// error interface and exposes the raw HTTP status code, the parsed
// `{"errors":[...]}` envelope, and the original response body for callers
// that want to inspect either form.
//
// Reference: https://developers.asana.com/docs/errors
type APIError struct {
	StatusCode int
	Errors     []APIErrorDetail
	RawBody    string
}

// APIErrorDetail is a single entry in Asana's `errors` array.
type APIErrorDetail struct {
	Message string `json:"message"`
	Help    string `json:"help,omitempty"`
	Phrase  string `json:"phrase,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.RawBody)
	}
	parts := make([]string, 0, len(e.Errors))
	for _, d := range e.Errors {
		parts = append(parts, d.Message)
	}
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, strings.Join(parts, "; "))
}

// FirstMessage returns the first error message from the envelope, or the raw
// body if the envelope was unparseable. Useful for compact UI rendering.
func (e *APIError) FirstMessage() string {
	if len(e.Errors) > 0 {
		return e.Errors[0].Message
	}
	return e.RawBody
}

// FirstHelp returns the help URL associated with the first error, if any.
func (e *APIError) FirstHelp() string {
	if len(e.Errors) > 0 {
		return e.Errors[0].Help
	}
	return ""
}

// IsNotFound reports whether the API responded with 404.
func (e *APIError) IsNotFound() bool { return e.StatusCode == 404 }

// IsUnauthorized reports whether the API responded with 401.
func (e *APIError) IsUnauthorized() bool { return e.StatusCode == 401 }

// IsForbidden reports whether the API responded with 403.
func (e *APIError) IsForbidden() bool { return e.StatusCode == 403 }

// IsRateLimited reports whether the API responded with 429.
func (e *APIError) IsRateLimited() bool { return e.StatusCode == 429 }

// AsAPIError returns the underlying *APIError if err is one. Returns
// (nil, false) for any other error type.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// parseAPIError builds an *APIError from a non-2xx HTTP response body.
// Falls back to a body-only APIError if the envelope can't be parsed.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		RawBody:    string(body),
	}
	var envelope struct {
		Errors []APIErrorDetail `json:"errors"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil {
		apiErr.Errors = envelope.Errors
	}
	return apiErr
}
