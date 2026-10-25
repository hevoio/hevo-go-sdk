package pipelines_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestGetPipelinePositionMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/position").ThenReply(http.StatusOK,
		`{"code": "Success", "data": [{"fully_qualified_name": "public.orders", "name": "orders", "position": "0/16B3748"}], "has_more": false}`)

	response, err := client.GetPipelinePosition().PipelineID(42).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].Name, "orders")
	testutils.AssertEqual(t, *response.Data[0].Position, "0/16B3748")
}
