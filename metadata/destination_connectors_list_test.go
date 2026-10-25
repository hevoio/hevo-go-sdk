package metadata_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestListDestinationConnectorsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/metadata/destinations/connectors").ThenReply(http.StatusOK, connectorsList)

	response, err := client.ListDestinationConnectors().Limit(500).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].ConnectorID, "hubspot")
	testutils.AssertEqual(t, response.Data[0].Description == nil, true)
}
