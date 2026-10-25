package httputils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// clientFunc lets a test answer requests with a function.
type clientFunc func(req *http.Request) (*http.Response, error)

func (f clientFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func reply(status int, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}

	return &http.Response{StatusCode: status, Header: header, Body: io.NopCloser(strings.NewReader(""))}
}

func TestRequestDoRetries429OnPostWithTheSameBody(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	req := Request{
		Method: http.MethodPost,
		URL:    "https://us.hevodata.com/api/v2/pipelines",
		Body:   []byte(`{"name":"p"}`),
		Client: clientFunc(func(r *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"name":"p"}` {
				t.Errorf("attempt %d body = %q", calls.Load()+1, body)
			}

			if calls.Add(1) < 3 {
				return reply(http.StatusTooManyRequests, http.Header{headerRetryAfter: {"0"}}), nil
			}

			return reply(http.StatusOK, nil), nil
		}),
		HandleRateLimits: true,
		MaxRetryAttempts: 5,
	}

	_, status, _, err := req.Do(t.Context())
	if err != nil || status != http.StatusOK || calls.Load() != 3 {
		t.Fatalf("status %d after %d calls (err %v), want 200 after 3", status, calls.Load(), err)
	}
}

func TestRequestDoRetriesGatewayErrorsOnlyWhenIdempotent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		method        string
		expectedCalls int32
	}{
		{name: "get retries", method: http.MethodGet, expectedCalls: 3},
		{name: "delete retries", method: http.MethodDelete, expectedCalls: 3},
		{name: "post does not", method: http.MethodPost, expectedCalls: 1},
		{name: "patch does not", method: http.MethodPatch, expectedCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var calls atomic.Int32

			req := Request{
				Method: tt.method,
				URL:    "https://us.hevodata.com/api/v2/pipelines/1",
				Client: clientFunc(func(*http.Request) (*http.Response, error) {
					calls.Add(1)

					return reply(http.StatusGatewayTimeout, http.Header{headerRetryAfter: {"0"}}), nil
				}),
				MaxRetryAttempts: 2,
			}
			req.backoffOverride = time.Millisecond

			if _, _, _, err := req.Do(t.Context()); err != nil {
				t.Fatal(err)
			}

			if got := calls.Load(); got != tt.expectedCalls {
				t.Errorf("calls = %d, want %d", got, tt.expectedCalls)
			}
		})
	}
}

func TestRequestDoLeaves429AloneWhenNotHandlingRateLimits(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	req := Request{
		Method: http.MethodGet,
		URL:    "https://us.hevodata.com/api/v2/pipelines",
		Client: clientFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)

			return reply(http.StatusTooManyRequests, http.Header{headerRetryAfter: {"0"}}), nil
		}),
		MaxRetryAttempts: 5,
	}

	_, status, _, _ := req.Do(t.Context())
	if status != http.StatusTooManyRequests || calls.Load() != 1 {
		t.Errorf("status %d after %d calls, want 429 after 1", status, calls.Load())
	}
}

func TestRequestDoStopsWaitingWhenContextEnds(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()

	req := Request{
		Method: http.MethodGet,
		URL:    "https://us.hevodata.com/api/v2/pipelines",
		Client: clientFunc(func(*http.Request) (*http.Response, error) {
			return reply(http.StatusTooManyRequests, http.Header{headerRetryAfter: {"3600"}}), nil
		}),
		HandleRateLimits: true,
		MaxRetryAttempts: 5,
	}

	if _, _, _, err := req.Do(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want DeadlineExceeded", err)
	}
}

func TestRateLimitWait(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		header     http.Header
		expected   time.Duration
		expectedOK bool
	}{
		{name: "retry-after seconds", header: http.Header{headerRetryAfter: {"7"}}, expected: 7 * time.Second, expectedOK: true},
		{name: "reset without retry-after", header: http.Header{"X-Ratelimit-Reset": {"12"}}, expected: 12 * time.Second, expectedOK: true},
		{
			name:       "retry-after wins over reset",
			header:     http.Header{headerRetryAfter: {"3"}, "X-Ratelimit-Reset": {"12"}},
			expected:   3 * time.Second,
			expectedOK: true,
		},
		{name: "past http date", header: http.Header{headerRetryAfter: {"Mon, 02 Jan 2006 15:04:05 GMT"}}, expectedOK: true},
		{name: "no headers", header: http.Header{}},
		{name: "unparseable", header: http.Header{headerRetryAfter: {"soon"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := rateLimitWait(tt.header)
			if got != tt.expected || ok != tt.expectedOK {
				t.Errorf("rateLimitWait() = %v, %v, want %v, %v", got, ok, tt.expected, tt.expectedOK)
			}
		})
	}
}

func TestBackoff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		attempt int
		floor   time.Duration
	}{
		{name: "first retry", attempt: 0, floor: time.Second},
		{name: "fourth retry", attempt: 3, floor: 8 * time.Second},
		{name: "capped", attempt: 40, floor: maxBackoff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := backoff(tt.attempt)
			if got < tt.floor || got > tt.floor+tt.floor/5 {
				t.Errorf("backoff(%d) = %v, want within 20%% above %v", tt.attempt, got, tt.floor)
			}
		})
	}
}

func bodyOf(s string) io.ReadCloser { return io.NopCloser(strings.NewReader(s)) }
