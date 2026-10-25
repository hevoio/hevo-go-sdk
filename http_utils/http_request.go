package httputils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxBackoff = 30 * time.Second

// Request is one call to the API, with its retry policy:
//
//   - 429 on any method: the gateway rejected the call before it reached Hevo. It waits for
//     Retry-After, else X-RateLimit-Reset, else backs off.
//   - 502, 503, 504 and network errors, on GET and DELETE only. A POST or PATCH that failed this
//     way may have been applied, so it goes back to the caller to reconcile.
type Request struct {
	Method  string
	URL     string
	Body    []byte
	Queries url.Values
	Headers map[string]string
	Client  HTTPClient

	HandleRateLimits bool
	MaxRetryAttempts int

	// backoffOverride replaces the exponential backoff in tests.
	backoffOverride time.Duration
}

// Do sends the request and returns the body and status of the last attempt.
func (req *Request) Do(ctx context.Context) ([]byte, int, http.Header, error) {
	if req.Client == nil {
		return nil, 0, nil, errors.New("hevo: http client is not provided")
	}

	for attempt := 0; ; attempt++ {
		body, status, header, err := req.send(ctx)

		wait, retry := req.shouldRetry(status, header, err, attempt)
		if !retry {
			return body, status, header, err
		}

		debugf("  - Retrying in %s (attempt %d)\n", wait, attempt+1)

		if err := contextDelay(ctx, wait); err != nil {
			return nil, 0, nil, err
		}
	}
}

func (req *Request) send(ctx context.Context) ([]byte, int, http.Header, error) {
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("hevo: building request: %w", err)
	}

	if len(req.Queries) > 0 {
		httpReq.URL.RawQuery = req.Queries.Encode()
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	req.debugRequest(httpReq)

	resp, err := req.Client.Do(httpReq)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("hevo: %s %s: %w", req.Method, httpReq.URL.Path, err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("hevo: reading response: %w", err)
	}

	debugf("  - HTTP Response:\n    - Status Code: %d\n    - Body: %s\n---\n", resp.StatusCode, body)

	return body, resp.StatusCode, resp.Header, nil
}

func (req *Request) shouldRetry(status int, header http.Header, err error, attempt int) (time.Duration, bool) {
	if attempt >= req.MaxRetryAttempts || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 0, false
	}

	wait := backoff(attempt)
	if req.backoffOverride > 0 {
		wait = req.backoffOverride
	}

	isIdempotent := req.Method == http.MethodGet || req.Method == http.MethodDelete
	if err != nil {
		return wait, isIdempotent
	}

	switch status {
	case http.StatusTooManyRequests:
		if !req.HandleRateLimits {
			return 0, false
		}

		if limitWait, ok := rateLimitWait(header); ok {
			return limitWait, true
		}

		return wait, true
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return wait, isIdempotent
	default:
		return 0, false
	}
}

func (req *Request) debugRequest(httpReq *http.Request) {
	if !debug.isEnabled {
		return
	}

	debugf("---\nDebug:\n  - HTTP Request:\n    - Method: %s\n    - URL: %s\n    - Body: %s\n    - Headers:\n",
		req.Method, httpReq.URL, req.Body)

	for k, v := range req.Headers {
		if k == "Authorization" && !debug.isAuthEnabled {
			v = "<omitted>"
		}

		debugf("      - %s: %s\n", k, v)
	}
}

// rateLimitWait reads how long to wait from Retry-After (seconds or an HTTP date), else from
// X-RateLimit-Reset (seconds until the gateway's window resets).
func rateLimitWait(h http.Header) (time.Duration, bool) {
	if v := h.Get(headerRetryAfter); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second, true
		}

		if at, err := http.ParseTime(v); err == nil {
			return max(time.Until(at), 0), true
		}
	}

	if secs, err := strconv.Atoi(h.Get(headerReset)); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second, true
	}

	return 0, false
}

// backoff is exponential from 1 s, capped at 30 s, plus up to 20% jitter.
func backoff(attempt int) time.Duration {
	d := min(time.Second<<min(attempt, 5), maxBackoff)

	return d + rand.N(d/5+1) //nolint:gosec // jitter needs no cryptographic randomness
}

func contextDelay(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("hevo: interrupted while waiting to retry: %w", ctx.Err())
	case <-t.C:
		return nil
	}
}
