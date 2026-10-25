package jobs_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestCancelPipelineJobMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPost, "/api/v2/pipelines/42/jobs/j-1/actions/cancel").ThenReply(http.StatusAccepted,
		`{"code": "Accepted", "data": {"id": "j-1", "status": "CANCELING"}}`)

	response, err := client.CancelPipelineJob().PipelineID(42).JobID("j-1").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data.Status, "CANCELING")
}
