package jobs_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestGetPipelineJobMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/jobs/j-1").ThenReply(http.StatusOK,
		`{"code": "Success", "data": `+jobJSON+`}`)

	response, err := client.GetPipelineJob().PipelineID(42).JobID("j-1").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data.Status, "COMPLETED")
}

func TestGetPipelineJobRequiresJobID(t *testing.T) {
	t.Parallel()

	client, _ := testutils.CreateTestClient()
	if _, err := client.GetPipelineJob().PipelineID(42).Do(t.Context()); err == nil {
		t.Fatal("want an error without a job id")
	}
}
