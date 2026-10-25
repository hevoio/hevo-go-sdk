package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/url"
)

// HTTPService is embedded in every operation's Service. It carries the client's settings, so an
// operation only names its method, path, body and response.
type HTTPService struct {
	BaseURL          string
	CommonHeaders    map[string]string
	Client           HTTPClient
	HandleRateLimits bool
	MaxRetryAttempts int
	// Slots caps the requests in flight across every Service of one client; nil means no cap.
	Slots chan struct{}
}

// Do sends method path with requestBody as JSON, and decodes a 2xx body into response. A non-2xx
// response is a *ResponseError.
func (s HTTPService) Do(
	ctx context.Context,
	method string,
	path string,
	requestBody any,
	queries url.Values,
	response any,
) error {
	var body []byte

	if requestBody != nil {
		var err error

		body, err = json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("hevo: encoding request body: %w", err)
		}
	}

	headers := make(map[string]string, len(s.CommonHeaders)+1)
	maps.Copy(headers, s.CommonHeaders)

	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	if err := s.acquire(ctx); err != nil {
		return err
	}
	defer s.release()

	r := Request{
		Method:           method,
		URL:              s.BaseURL + path,
		Body:             body,
		Queries:          queries,
		Headers:          headers,
		Client:           s.Client,
		HandleRateLimits: s.HandleRateLimits,
		MaxRetryAttempts: s.MaxRetryAttempts,
	}

	respBody, status, header, err := r.Do(ctx)
	if err != nil {
		return err
	}

	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return newResponseError(status, header, respBody)
	}

	if response == nil || len(respBody) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBody, response); err != nil {
		return fmt.Errorf("hevo: decoding %s %s response: %w", method, path, err)
	}

	return nil
}

func (s HTTPService) acquire(ctx context.Context) error {
	if s.Slots == nil {
		return nil
	}

	select {
	case s.Slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("hevo: waiting for a request slot: %w", ctx.Err())
	}
}

func (s HTTPService) release() {
	if s.Slots != nil {
		<-s.Slots
	}
}
