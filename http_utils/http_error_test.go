package httputils

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewResponseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		status          int
		body            string
		expectedMessage string
		expectedString  string
	}{
		{
			name:            "api error body",
			status:          http.StatusNotFound,
			body:            `{"code":"RESOURCE_NOT_FOUND","message":"No Pipeline 7."}`,
			expectedMessage: "No Pipeline 7.",
			expectedString:  "hevo: 404 RESOURCE_NOT_FOUND: No Pipeline 7.",
		},
		{
			name:            "non-json body",
			status:          http.StatusBadGateway,
			body:            "<html>bad gateway</html>",
			expectedMessage: "<html>bad gateway</html>",
			expectedString:  "hevo: 502: <html>bad gateway</html>",
		},
		{name: "empty body", status: http.StatusConflict, expectedString: "hevo: 409: Conflict"},
		{
			name:            "error reference",
			status:          http.StatusInternalServerError,
			body:            `{"code":"INTERNAL_ERROR","message":"Something went wrong.","error_reference_id":"abc"}`,
			expectedMessage: "Something went wrong.",
			expectedString:  "hevo: 500 INTERNAL_ERROR: Something went wrong. (error reference abc)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := newResponseError(tt.status, http.Header{}, []byte(tt.body))
			if err.Message != tt.expectedMessage {
				t.Errorf("Message = %q, want %q", err.Message, tt.expectedMessage)
			}

			if got := err.Error(); got != tt.expectedString {
				t.Errorf("Error() = %q, want %q", got, tt.expectedString)
			}
		})
	}
}

func TestNewResponseErrorGateway429(t *testing.T) {
	t.Parallel()

	h := http.Header{}
	h.Set(headerLimit, "60")
	h.Set(headerRemaining, "0")
	h.Set(headerReset, "42")
	h.Set(headerBucket, "write")
	h.Set(headerRetryAfter, "42")

	err := newResponseError(http.StatusTooManyRequests, h, nil)
	if err.Code != CodeRateLimited || !IsRateLimited(err) {
		t.Fatalf("err = %v, want RATE_LIMITED", err)
	}

	expected := RateLimit{Limit: 60, Remaining: 0, Reset: 42 * time.Second, Bucket: "write", RetryAfter: 42 * time.Second}
	if *err.RateLimit != expected {
		t.Errorf("RateLimit = %+v, want %+v", *err.RateLimit, expected)
	}

	const expectedString = `hevo: 429 RATE_LIMITED: too many requests for this team's API quota` +
		` (limit 60 per window on "write"; retry after 42s)`
	if got := err.Error(); got != expectedString {
		t.Errorf("Error() = %q, want %q", got, expectedString)
	}
}

func TestParseRateLimit(t *testing.T) {
	t.Parallel()

	if got := ParseRateLimit(http.Header{}); got != nil {
		t.Errorf("ParseRateLimit() = %+v, want nil without %s", got, headerLimit)
	}

	h := http.Header{}
	h.Set(headerLimit, "100")
	h.Set(headerRemaining, "99")
	h.Set(headerReset, "60")

	expected := RateLimit{Limit: 100, Remaining: 99, Reset: time.Minute}
	if got := ParseRateLimit(h); got == nil || *got != expected {
		t.Errorf("ParseRateLimit() = %+v, want %+v", got, expected)
	}
}

func TestStatusPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		status    int
		predicate func(error) bool
	}{
		{name: "not found", status: http.StatusNotFound, predicate: IsNotFound},
		{name: "conflict", status: http.StatusConflict, predicate: IsConflict},
		{name: "rate limited", status: http.StatusTooManyRequests, predicate: IsRateLimited},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := fmt.Errorf("reading pipeline: %w", newResponseError(tt.status, http.Header{}, nil))
			if !tt.predicate(err) {
				t.Errorf("predicate(%v) = false, want true through wrapping", err)
			}

			if tt.predicate(newResponseError(http.StatusBadRequest, http.Header{}, nil)) || tt.predicate(errors.New("x")) {
				t.Error("predicate matched a 400 or a plain error")
			}
		})
	}
}
