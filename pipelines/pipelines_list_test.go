package pipelines_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestListPipelinesMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			testutils.AssertEqual(t, q["state"], []string{"Initialized", "Disabled"})
			testutils.AssertEqual(t, q.Get("connector_id"), "hubspot")
			testutils.AssertEqual(t, q.Get("destination_id"), "7")
			testutils.AssertEqual(t, q.Get("limit"), "20")

			body := `{"code": "Success", "data": [` + pipelineJSON + `], "next_cursor": null, "has_more": false}`

			return mock.NewResponse(req, http.StatusOK, body), nil
		})

	response, err := client.ListPipelines().
		States("Initialized", "Disabled").
		ConnectorID("hubspot").
		DestinationID(7).
		Limit(20).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, len(response.Data), 1)
	assertPipeline(t, response.Data[0])
	testutils.AssertEqual(t, response.HasMore, false)
}
