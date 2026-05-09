package asana

import (
	"errors"
	"strings"
	"testing"
)

func TestParseAPIErrorWithEnvelope(t *testing.T) {
	body := []byte(`{"errors":[{"message":"task: Not a recognized ID: 999","help":"https://developers.asana.com/docs/errors"}]}`)
	apiErr := parseAPIError(404, body)
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if len(apiErr.Errors) != 1 {
		t.Fatalf("Errors len = %d, want 1", len(apiErr.Errors))
	}
	if apiErr.Errors[0].Message != "task: Not a recognized ID: 999" {
		t.Errorf("Errors[0].Message = %q", apiErr.Errors[0].Message)
	}
	if apiErr.Errors[0].Help != "https://developers.asana.com/docs/errors" {
		t.Errorf("Errors[0].Help = %q", apiErr.Errors[0].Help)
	}
}

func TestParseAPIErrorMultiple(t *testing.T) {
	body := []byte(`{"errors":[{"message":"first"},{"message":"second","help":"h"}]}`)
	apiErr := parseAPIError(400, body)
	if len(apiErr.Errors) != 2 {
		t.Fatalf("Errors len = %d, want 2", len(apiErr.Errors))
	}
	if !strings.Contains(apiErr.Error(), "first; second") {
		t.Errorf("Error() = %q, want it to join messages with semicolons", apiErr.Error())
	}
}

func TestParseAPIErrorMalformedFallsBackToBody(t *testing.T) {
	body := []byte(`<html>500 Internal Server Error</html>`)
	apiErr := parseAPIError(500, body)
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d", apiErr.StatusCode)
	}
	if len(apiErr.Errors) != 0 {
		t.Errorf("Errors should be empty for unparseable body, got %d", len(apiErr.Errors))
	}
	if apiErr.RawBody != string(body) {
		t.Errorf("RawBody = %q", apiErr.RawBody)
	}
	got := apiErr.Error()
	if !strings.Contains(got, "API error (500)") || !strings.Contains(got, "Internal Server Error") {
		t.Errorf("Error() = %q", got)
	}
}

func TestAPIErrorStatusHelpers(t *testing.T) {
	cases := []struct {
		name string
		code int
		fn   func(*APIError) bool
		want bool
	}{
		{"401 unauthorized", 401, (*APIError).IsUnauthorized, true},
		{"403 forbidden", 403, (*APIError).IsForbidden, true},
		{"404 not found", 404, (*APIError).IsNotFound, true},
		{"429 rate limited", 429, (*APIError).IsRateLimited, true},
		{"200 not unauthorized", 200, (*APIError).IsUnauthorized, false},
		{"500 not unauthorized", 500, (*APIError).IsUnauthorized, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := &APIError{StatusCode: tc.code}
			if got := tc.fn(err); got != tc.want {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestAsAPIErrorRoundTrip(t *testing.T) {
	original := &APIError{StatusCode: 404, Errors: []APIErrorDetail{{Message: "x"}}}
	var asError error = original

	got, ok := AsAPIError(asError)
	if !ok {
		t.Fatal("AsAPIError should return true for *APIError")
	}
	if got.StatusCode != 404 {
		t.Errorf("StatusCode = %d", got.StatusCode)
	}

	plain := errors.New("not an api error")
	if _, ok := AsAPIError(plain); ok {
		t.Error("AsAPIError should return false for plain errors")
	}
}

func TestFirstMessageFallsBackToRawBody(t *testing.T) {
	apiErr := &APIError{StatusCode: 500, RawBody: "boom"}
	if apiErr.FirstMessage() != "boom" {
		t.Errorf("FirstMessage = %q", apiErr.FirstMessage())
	}

	apiErr.Errors = []APIErrorDetail{{Message: "structured"}}
	if apiErr.FirstMessage() != "structured" {
		t.Errorf("FirstMessage = %q", apiErr.FirstMessage())
	}
}
