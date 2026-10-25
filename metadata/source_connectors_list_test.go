package metadata_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestListSourceConnectorsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/metadata/sources/connectors").ThenReply(http.StatusOK, connectorsList)

	response, err := client.ListSourceConnectors().Limit(500).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data[0].ConnectorID, "hubspot")
	testutils.AssertEqual(t, response.Data[0].Description == nil, true)
}

func TestListSourceConnectorsFiltersByType(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodGet, "/api/v2/metadata/sources/connectors").ThenReply(http.StatusOK, connectorsList)

	if _, err := client.ListSourceConnectors().Type("HUBSPOT").Do(t.Context()); err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, mockClient.Interactions()[0].Request.URL.Query().Get("type"), "HUBSPOT")
}
