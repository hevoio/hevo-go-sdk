package pipelines_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestGetPipelineMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			return mock.NewResponse(req, http.StatusOK, pipelineResponse("Success")), nil
		})

	response, err := client.GetPipeline().PipelineID(42).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Success")
	assertPipeline(t, response.Data)
}

func TestGetPipelineRequiresPipelineID(t *testing.T) {
	t.Parallel()

	client, _ := testutils.CreateTestClient()
	if _, err := client.GetPipeline().Do(t.Context()); err == nil {
		t.Fatal("want an error without a pipeline id")
	}
}
