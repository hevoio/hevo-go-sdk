package pipelines_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestSyncPipelineNowMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPost, "/api/v2/pipelines/42/actions/sync-now").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			return mock.NewResponse(req, http.StatusAccepted, pipelineResponse("Accepted")), nil
		})

	response, err := client.SyncPipelineNow().PipelineID(42).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Accepted")
	assertPipeline(t, response.Data)
}

func TestSyncPipelineNowRequiresPipelineID(t *testing.T) {
	t.Parallel()

	client, _ := testutils.CreateTestClient()
	if _, err := client.SyncPipelineNow().Do(t.Context()); err == nil {
		t.Fatal("want an error without a pipeline id")
	}
}
