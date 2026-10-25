package destinations_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestListDestinationsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/destinations").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, req.URL.Query().Get("limit"), "10")
			testutils.AssertEqual(t, req.URL.Query().Get("cursor"), "MTIzNDU=")

			body := `{"code": "Success", "data": [` + destinationJSON + `], "next_cursor": "Njc4OQ==", "has_more": true}`

			return mock.NewResponse(req, http.StatusOK, body), nil
		})

	response, err := client.ListDestinations().Limit(10).Cursor("MTIzNDU=").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, len(response.Data), 1)
	assertDestination(t, response.Data[0])
	testutils.AssertEqual(t, response.NextCursor, "Njc4OQ==")
	testutils.AssertEqual(t, response.HasMore, true)
}
