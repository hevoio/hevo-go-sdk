package httputils

import (
	"context"
	"errors"
	"net/http"
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
