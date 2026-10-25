package jobs_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestListPipelineJobsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/jobs").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, req.URL.Query()["status"], []string{"COMPLETED"})
			testutils.AssertEqual(t, req.URL.Query()["type"], []string{"INCREMENTAL", "HISTORICAL"})

			return mock.NewResponse(req, http.StatusOK, `{"code": "Success", "data": [`+jobJSON+`], "has_more": false}`), nil
		})

	response, err := client.ListPipelineJobs().
		PipelineID(42).
		Statuses("COMPLETED").
		Types("INCREMENTAL", "HISTORICAL").
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].EventsLoaded, int64(9))
	testutils.AssertEqual(t, response.Data[0].Latency.Mean, int64(400))
}
