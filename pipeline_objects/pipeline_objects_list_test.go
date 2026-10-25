package pipelineobjects_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestListPipelineObjectsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/objects").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, req.URL.Query()["status"], []string{"ACTIVE", "INACTIVE"})

			body := `{"code": "Success", "data": [` + objectJSON + `], "next_cursor": "c2", "has_more": true}`

			return mock.NewResponse(req, http.StatusOK, body), nil
		})

	response, err := client.ListPipelineObjects().PipelineID(42).Statuses("ACTIVE", "INACTIVE").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].Name, "public.orders")
	testutils.AssertEqual(t, response.Data[0].Destination.Table, "hubspot_orders")
	testutils.AssertEqual(t, response.Data[0].Fields[0].PrimaryKey, true)
	testutils.AssertEqual(t, response.NextCursor, "c2")
}
