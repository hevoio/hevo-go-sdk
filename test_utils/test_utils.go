// Package testutils holds what the unit and end-to-end tests share.
package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/hevoio/hevo-go"
	"github.com/hevoio/hevo-go/tests/mock"
)

const (
	testAccessKey = "test_access_key"
	testSecretKey = "test_secret_key"
)

// CreateTestClient returns a client that sends every request to a new mock, without retries
// so a test sees each response as it is.
func CreateTestClient() (*hevo.Client, *mock.HTTPClient) {
	client := hevo.New(testAccessKey, testSecretKey)
	mockClient := mock.NewHTTPClient()
	client.SetHTTPClient(mockClient)
	client.SetMaxRetryAttempts(0)

	return client, mockClient
}

// RequestBodyToJSON decodes the request body, and puts it back so the request can be read again.
func RequestBodyToJSON(t *testing.T, req *http.Request) map[string]any {
	t.Helper()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("reading request body: %v", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	result := map[string]any{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("decoding request body %q: %v", body, err)
	}

	return result
}

// JSONOf re-encodes value the way a request body is decoded, so maps of any shape compare equal.
func JSONOf(t *testing.T, value any) any {
	t.Helper()

	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}

	return decoded
}

// AssertEqual fails the test when actual and expected differ.
func AssertEqual(t *testing.T, actual, expected any) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %#v, want %#v", actual, expected)
	}
}

// E2EClient returns a client for the account in HEVO_ACCESS_KEY, HEVO_SECRET_KEY and
// HEVO_BASE_URL, and skips the test when they are not set.
func E2EClient(t *testing.T) *hevo.Client {
	t.Helper()

	accessKey, secretKey, baseURL := os.Getenv("HEVO_ACCESS_KEY"), os.Getenv("HEVO_SECRET_KEY"), os.Getenv("HEVO_BASE_URL")
	if accessKey == "" || secretKey == "" || baseURL == "" {
		t.Skip("set HEVO_ACCESS_KEY, HEVO_SECRET_KEY and HEVO_BASE_URL to run end-to-end tests")
	}

	client := hevo.New(accessKey, secretKey)
	client.BaseURL(baseURL)

	return client
}
