package destinations_test

import (
	"net/http"
	"testing"

	"github.com/hevoio/hevo-go/destinations"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestCreateDestinationMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPost, "/api/v2/destinations").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"name":         "Snowflake",
				"connector_id": "snowflake",
				"type":         "SNOWFLAKE",
				"config":       map[string]any{"account_url": "acme.snowflakecomputing.com", "password": "secret"},
			})

			return mock.NewResponse(req, http.StatusOK, destinationResponse), nil
		})

	response, err := client.CreateDestination().
		Name("Snowflake").
		ConnectorID("snowflake").
		Type("SNOWFLAKE").
		Config(map[string]any{"account_url": "acme.snowflakecomputing.com", "password": "secret"}).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	assertDestination(t, response.Data)
}

func assertDestination(t *testing.T, d destinations.DestinationDetails) {
	t.Helper()

	testutils.AssertEqual(t, d.ID, int64(7))
	testutils.AssertEqual(t, d.Name, "Snowflake")
	testutils.AssertEqual(t, d.ConnectorID, "snowflake")
	testutils.AssertEqual(t, d.Type, "SNOWFLAKE")
	testutils.AssertEqual(t, d.Status, "CREATED")
	testutils.AssertEqual(t, d.Config["password"], "******")
	testutils.AssertEqual(t, d.CreatedAt, int64(1760000000000))
	testutils.AssertEqual(t, d.UpdatedBy, "ops@acme.com")
}

func TestCreateDestinationSendsAnEmptyConfig(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodPost, "/api/v2/destinations").ThenCall(func(req *http.Request) (*http.Response, error) {
		testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
			"name": "d", "connector_id": "blackhole_sf", "config": map[string]any{},
		})

		return mock.NewResponse(req, http.StatusOK, destinationResponse), nil
	})

	if _, err := client.CreateDestination().Name("d").ConnectorID("blackhole_sf").Config(map[string]any{}).Do(t.Context()); err != nil {
		t.Fatal(err)
	}
}
