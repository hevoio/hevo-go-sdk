package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Error codes the Hevo API returns in an error body. Each always comes with the same status.
const (
	CodeInvalidInput     = "INVALID_INPUT"
	CodeNotAuthorized    = "NOT_AUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	CodeResourceConflict = "RESOURCE_CONFLICT"
	CodeRateLimited      = "RATE_LIMITED"
	CodeInternalError    = "INTERNAL_ERROR"
)

const (
	headerRetryAfter = "Retry-After"
	headerLimit      = "X-RateLimit-Limit"
	headerRemaining  = "X-RateLimit-Remaining"
	headerReset      = "X-RateLimit-Reset"
	headerBucket     = "X-RateLimit-Bucket"
)

// ResponseError is a response with a non-2xx status.
type ResponseError struct {
	StatusCode       int    `json:"-"`
	Code             string `json:"code"`
	Message          string `json:"message"`
	ErrorReferenceID string `json:"error_reference_id"`
	APIDocURL        string `json:"api_doc_url"`
	// RateLimit is the team's quota for the bucket the call counted against, when the gateway
	// reported it. On a 429, its RetryAfter says when to try again.
	RateLimit *RateLimit `json:"-"`
}

// RateLimit is the quota the gateway reports on every API response.
type RateLimit struct {
	Limit      int
	Remaining  int
	Reset      time.Duration
	Bucket     string
	RetryAfter time.Duration
}

func (e *ResponseError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}

	s := "hevo: " + strconv.Itoa(e.StatusCode)
	if e.Code != "" {
		s += " " + e.Code
	}

	s += ": " + msg
	if e.ErrorReferenceID != "" {
		s += " (error reference " + e.ErrorReferenceID + ")"
	}

	if e.StatusCode == http.StatusTooManyRequests && e.RateLimit != nil {
		s += fmt.Sprintf(" (limit %d per window on %q; retry after %s)",
			e.RateLimit.Limit, e.RateLimit.Bucket, e.RateLimit.RetryAfter)
	}

	return s
}

func newResponseError(status int, header http.Header, body []byte) *ResponseError {
	e := &ResponseError{}
	if len(body) > 0 && json.Unmarshal(body, e) != nil {
		e.Message = string(body)
	}

	e.StatusCode = status
	e.RateLimit = ParseRateLimit(header)

	if status == http.StatusTooManyRequests && e.Code == "" {
		// The gateway rejects an over-quota call with an empty body.
		e.Code = CodeRateLimited
		e.Message = "too many requests for this team's API quota"
	}

	return e
}

// ParseRateLimit reads the gateway's X-RateLimit-* and Retry-After headers. It returns nil when
// the response carries none, for example on a route without a quota.
func ParseRateLimit(h http.Header) *RateLimit {
	limit, ok := headerInt(h, headerLimit)
	if !ok {
		return nil
	}

	rl := &RateLimit{Limit: limit, Bucket: h.Get(headerBucket)}
	rl.Remaining, _ = headerInt(h, headerRemaining)

	if reset, ok := headerInt(h, headerReset); ok {
		rl.Reset = time.Duration(reset) * time.Second
	}

	if h.Get(headerRetryAfter) != "" {
		rl.RetryAfter, _ = rateLimitWait(h)
	}

	return rl
}

// IsNotFound reports whether err is a 404, for example a resource deleted outside the caller.
func IsNotFound(err error) bool { return hasStatus(err, http.StatusNotFound) }

// IsConflict reports whether err is a 409: the Pipeline is mid-transition or has an edit open.
func IsConflict(err error) bool { return hasStatus(err, http.StatusConflict) }

// IsRateLimited reports whether err is a 429 that outlasted the client's retries.
func IsRateLimited(err error) bool { return hasStatus(err, http.StatusTooManyRequests) }

func hasStatus(err error, status int) bool {
	e, ok := errors.AsType[*ResponseError](err)

	return ok && e.StatusCode == status
}

func headerInt(h http.Header, name string) (int, bool) {
	v, err := strconv.Atoi(h.Get(name))

	return v, err == nil
}
