package httputils

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPServiceDoDecodesAndSetsHeaders(t *testing.T) {
	t.Parallel()

	s := HTTPService{
		BaseURL:       "https://us.hevodata.com",
		CommonHeaders: map[string]string{"Authorization": "Basic x", "User-Agent": "hevo-go/test"},
		Client: clientFunc(func(r *http.Request) (*http.Response, error) {
			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q", got)
			}

			if got := r.URL.String(); got != "https://us.hevodata.com/api/v2/destinations?limit=5" {
				t.Errorf("URL = %q", got)
			}

			resp := reply(http.StatusOK, nil)
			resp.Body = bodyOf(`{"code":"Success","message":"done"}`)

			return resp, nil
		}),
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	limit := 5

	err := s.Do(t.Context(), http.MethodPost, "/api/v2/destinations", map[string]string{"name": "d"},
		PageQueries(&limit, nil), &response)
	if err != nil || response.Message != "done" {
		t.Fatalf("Do() = %v, response %+v", err, response)
	}
}

func TestHTTPServiceDoReturnsResponseError(t *testing.T) {
	t.Parallel()

	s := HTTPService{
		BaseURL: "https://us.hevodata.com",
		Client: clientFunc(func(*http.Request) (*http.Response, error) {
			resp := reply(http.StatusConflict, nil)
			resp.Body = bodyOf(`{"code":"RESOURCE_CONFLICT","message":"Edit in progress."}`)

			return resp, nil
		}),
	}

	err := s.Do(t.Context(), http.MethodPatch, "/api/v2/pipelines/1", map[string]string{}, nil, nil)
	if !IsConflict(err) {
		t.Fatalf("Do() = %v, want a 409", err)
	}
}

func TestHTTPServiceDoWaitsForASlot(t *testing.T) {
	t.Parallel()

	slots := make(chan struct{}, 1)
	slots <- struct{}{} // the only slot is taken

	s := HTTPService{BaseURL: "https://us.hevodata.com", Client: clientFunc(nil), Slots: slots}

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()

	if err := s.Do(ctx, http.MethodGet, "/api/v2/pipelines", nil, nil, nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Do() = %v, want DeadlineExceeded while no slot is free", err)
	}
}

func TestHTTPServiceDoFreesTheSlotWhileWaitingToRetry(t *testing.T) {
	t.Parallel()

	limited := make(chan struct{})

	s := HTTPService{
		BaseURL: "https://us.hevodata.com",
		Client: clientFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/api/v2/pipelines/1" {
				close(limited)

				return reply(http.StatusTooManyRequests, http.Header{headerRetryAfter: {"1"}}), nil
			}

			return reply(http.StatusOK, nil), nil
		}),
		HandleRateLimits: true,
		MaxRetryAttempts: 1,
		Slots:            make(chan struct{}, 1),
	}

	waitingCtx, stopWaiting := context.WithCancel(t.Context())
	waitingErr := make(chan error, 1)

	go func() {
		waitingErr <- s.Do(waitingCtx, http.MethodGet, "/api/v2/pipelines/1", nil, nil, nil)
	}()

	<-limited

	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()

	if err := s.Do(ctx, http.MethodGet, "/api/v2/pipelines", nil, nil, nil); err != nil {
		t.Errorf("Do() = %v while another call waits on a 429, want a free slot", err)
	}

	stopWaiting()

	if err := <-waitingErr; !errors.Is(err, context.Canceled) {
		t.Errorf("waiting Do() = %v, want Canceled", err)
	}
}

func TestHTTPServiceDoCapsTheRateLimitWait(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		header        http.Header
		expectedCalls int32
		expectedErr   bool
	}{
		{name: "retry-after over the cap", header: http.Header{headerRetryAfter: {"3600"}}, expectedCalls: 1, expectedErr: true},
		{name: "reset over the cap", header: http.Header{"X-Ratelimit-Reset": {"3600"}}, expectedCalls: 1, expectedErr: true},
		{name: "retry-after within the cap", header: http.Header{headerRetryAfter: {"0"}}, expectedCalls: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var calls atomic.Int32

			s := HTTPService{
				BaseURL: "https://us.hevodata.com",
				Client: clientFunc(func(*http.Request) (*http.Response, error) {
					if calls.Add(1) == 1 {
						header := tt.header.Clone()
						header.Set(headerLimit, "60")

						return reply(http.StatusTooManyRequests, header), nil
					}

					return reply(http.StatusOK, nil), nil
				}),
				HandleRateLimits: true,
				MaxRetryAttempts: 5,
			}

			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()

			err := s.Do(ctx, http.MethodGet, "/api/v2/pipelines", nil, nil, nil)

			if got := calls.Load(); got != tt.expectedCalls {
				t.Errorf("calls = %d, want %d", got, tt.expectedCalls)
			}

			if !tt.expectedErr {
				if err != nil {
					t.Errorf("Do() = %v, want nil", err)
				}

				return
			}

			if e, ok := errors.AsType[*ResponseError](err); !ok || !IsRateLimited(err) || e.RateLimit == nil {
				t.Errorf("Do() = %v, want a 429 *ResponseError with its RateLimit", err)
			}
		})
	}
}
