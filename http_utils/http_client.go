// Package httputils sends the SDK's requests to the Hevo API: authentication headers, JSON
// bodies, retries on rate limits and gateway errors, and error decoding.
package httputils

import "net/http"

// HTTPClient performs an HTTP request. *http.Client satisfies it; tests pass a mock.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
