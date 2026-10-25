package destinations_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestUpdateDestinationMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPatch, "/api/v2/destinations/7").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"name":   "Snowflake",
				"config": map[string]any{"warehouse": "LOAD_WH"},
			})

			return mock.NewResponse(req, http.StatusOK, destinationResponse), nil
		})

	response, err := client.UpdateDestination().
		DestinationID(7).
		Name("Snowflake").
		Config(map[string]any{"warehouse": "LOAD_WH"}).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	assertDestination(t, response.Data)
}
