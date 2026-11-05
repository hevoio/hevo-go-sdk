// Package mock is an in-memory HTTP client for unit tests. Stub a method and path, then assert
// on what was sent:
//
//	client := mock.NewHTTPClient()
//	handler := client.When(http.MethodGet, "/api/v2/pipelines/42").ThenReply(200, `{"code":"Success"}`)
//	// ... run the code under test with client ...
//	if handler.Interactions != 1 { t.Fatal("not called") }
//
// ThenCall builds the response from the request instead, to assert on the body sent.
package mock

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// HandlerFunc builds the response to one request.
type HandlerFunc func(req *http.Request) (*http.Response, error)

// Handler answers one stubbed method and path and counts its calls.
type Handler struct {
	reply        HandlerFunc
	Interactions int
}

// Interaction is one request the client received.
type Interaction struct {
	Handler *Handler
	Request *http.Request
}

type route struct {
	method string
	path   string
}

// HTTPClient answers requests from stubs. It is safe for concurrent use.
type HTTPClient struct {
	mu           sync.Mutex
	handlers     map[route]*Handler
	interactions []Interaction
}

// NewHTTPClient returns a client with no stubs.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{handlers: map[route]*Handler{}}
}

// When stubs method and path; the path excludes the host and query.
func (c *HTTPClient) When(method, path string) *Handler {
	c.mu.Lock()
	defer c.mu.Unlock()

	h := &Handler{}
	c.handlers[route{method: method, path: path}] = h

	return h
}

// ThenReply answers with status and body.
func (h *Handler) ThenReply(status int, body string) *Handler {
	h.reply = func(req *http.Request) (*http.Response, error) {
		return NewResponse(req, status, body), nil
	}

	return h
}

// ThenCall answers with what fn returns.
func (h *Handler) ThenCall(fn HandlerFunc) *Handler {
	h.reply = fn

	return h
}

// Do answers req from the stub for its method and path.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()

	h, ok := c.handlers[route{method: req.Method, path: req.URL.Path}]
	if ok {
		h.Interactions++
		c.interactions = append(c.interactions, Interaction{Handler: h, Request: req})
	}
	c.mu.Unlock()

	if !ok || h.reply == nil {
		return nil, fmt.Errorf("mock: no stub for %s %s", req.Method, req.URL.Path)
	}

	return h.reply(req)
}

// Interactions lists the requests received, in order.
func (c *HTTPClient) Interactions() []Interaction {
	c.mu.Lock()
	defer c.mu.Unlock()

	return append([]Interaction(nil), c.interactions...)
}

// NewResponse builds a JSON response to req.
func NewResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{
		Status:        http.StatusText(status),
		StatusCode:    status,
		Header:        http.Header{"Content-Type": {"application/json"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
	}
}
