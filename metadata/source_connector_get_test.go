package metadata_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestGetSourceConnectorMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/metadata/sources/connectors/hubspot").ThenReply(http.StatusOK, connectorDetails)

	response, err := client.GetSourceConnector().ConnectorID("hubspot").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data.Schema["type"], "object")
}
