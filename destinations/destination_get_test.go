package destinations_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestGetDestinationMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/destinations/7").ThenReply(http.StatusOK, destinationResponse)

	response, err := client.GetDestination().DestinationID(7).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Success")
	assertDestination(t, response.Data)
}

func TestGetDestinationRequiresID(t *testing.T) {
	t.Parallel()

	client, _ := testutils.CreateTestClient()
	if _, err := client.GetDestination().Do(t.Context()); err == nil {
		t.Fatal("want an error without a destination id")
	}
}

func TestGetDestinationNotFound(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodGet, "/api/v2/destinations/8").ThenReply(http.StatusNotFound,
		`{"code": "RESOURCE_NOT_FOUND", "message": "No Destination exists with id 8."}`)

	_, err := client.GetDestination().DestinationID(8).Do(t.Context())
	testutils.AssertEqual(t, err.Error(), "hevo: 404 RESOURCE_NOT_FOUND: No Destination exists with id 8.")
}
