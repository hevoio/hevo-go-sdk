package hevo_test

import (
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/hevoio/hevo-go"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestClientRegion(t *testing.T) {
	t.Parallel()

	client := hevo.New("key", "secret")
	if err := client.Region(" EU2 "); err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, client.HTTPService().BaseURL, "https://eu2.hevodata.com")

	if err := client.Region("mars"); err == nil {
		t.Error("want an error for an unknown region")
	}
}

func TestClientBaseURL(t *testing.T) {
	t.Parallel()

	client := hevo.New("key", "secret")
	client.BaseURL("http://localhost:9080/")
	testutils.AssertEqual(t, client.HTTPService().BaseURL, "http://localhost:9080")
}

func TestRegions(t *testing.T) {
	t.Parallel()

	expected := []string{"asia", "au", "eu", "eu2", "in", "us", "us2"}
	if got := hevo.Regions(); !slices.Equal(got, expected) {
		t.Errorf("Regions() = %v, want %v", got, expected)
	}
}

func TestClientSendsAuthAndUserAgent(t *testing.T) {
	t.Parallel()

	client := hevo.New("key", "secret")
	client.CustomUserAgent("terraform-provider-hevo/1.0.0")

	mockClient := mock.NewHTTPClient()
	client.SetHTTPClient(mockClient)
	mockClient.When(http.MethodGet, "/api/v2/pipelines/1").ThenCall(func(req *http.Request) (*http.Response, error) {
		user, pass, ok := req.BasicAuth()
		testutils.AssertEqual(t, []any{user, pass, ok}, []any{"key", "secret", true})

		ua := req.Header.Get("User-Agent")
		if !strings.HasPrefix(ua, hevo.DefaultUserAgent+" terraform-provider-hevo/1.0.0") {
			t.Errorf("User-Agent = %q", ua)
		}

		return mock.NewResponse(req, http.StatusOK, `{"code":"Success","data":{"id":1}}`), nil
	})

	if _, err := client.GetPipeline().PipelineID(1).Do(t.Context()); err != nil {
		t.Fatal(err)
	}
}
