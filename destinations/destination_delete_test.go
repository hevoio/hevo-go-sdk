package destinations_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestDeleteDestinationMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodDelete, "/api/v2/destinations/7").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, req.URL.Query().Get("confirm_override"), "true")

			return mock.NewResponse(req, http.StatusAccepted, `{"code": "Accepted", "data": `+destinationJSON+`}`), nil
		})

	response, err := client.DeleteDestination().DestinationID(7).ConfirmOverride(true).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Accepted")
}
