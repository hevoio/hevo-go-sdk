package jobs_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestListPipelineJobObjectsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/jobs/j-1/objects").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, req.URL.Query().Get("limit"), "5")

			body := `{"code": "Success", "data": [{"id": "o-1", "object_name": "public.orders", "status": "COMPLETED", "events_loaded": 9, "offset": null}], "has_more": false}`

			return mock.NewResponse(req, http.StatusOK, body), nil
		})

	response, err := client.ListPipelineJobObjects().PipelineID(42).JobID("j-1").Limit(5).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].ObjectName, "public.orders")
	testutils.AssertEqual(t, *response.Data[0].EventsLoaded, int64(9))
	testutils.AssertEqual(t, response.Data[0].EventsIngested == nil, true)
}
