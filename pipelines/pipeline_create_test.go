package pipelines_test

import (
	"net/http"
	"testing"

	"github.com/hevoio/hevo-go"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestCreatePipelineMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPost, "/api/v2/pipelines").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"name": "HubSpot to Snowflake",
				"source": map[string]any{
					"name":         "HubSpot",
					"connector_id": "hubspot",
					"type":         "HUBSPOT",
					"config":       map[string]any{"authentication_method": "OAUTH"},
					"auth": map[string]any{
						"client_access": map[string]any{"client_id": "cid", "client_secret": "csecret"},
						"refresh_token": "rtoken",
					},
				},
				"destination_id":     float64(7),
				"destination_prefix": "hubspot_",
				"replication_type":   "HISTORICAL_AND_INCREMENTAL",
				"load_mode":          "MERGE",
				"schema_evolution":   "ALLOW_ALL",
				"schedule":           map[string]any{"sync_type": "SCHEDULED", "frequency_minutes": float64(60)},
				"latency_alert":      map[string]any{"enabled": true, "threshold_minutes": float64(120)},
				"failure_handling_policy": map[string]any{
					"object_failure_level":     "PARTIAL",
					"object_failure_threshold": float64(5),
				},
				"options": map[string]any{"create_disabled": true, "include_all_objects": false},
			})

			return mock.NewResponse(req, http.StatusOK, pipelineResponse("Success")), nil
		})

	response, err := client.CreatePipeline().
		Name("HubSpot to Snowflake").
		Source(hevo.NewSource().
			Name("HubSpot").
			ConnectorID("hubspot").
			Type("HUBSPOT").
			Config(map[string]any{"authentication_method": "OAUTH"}).
			Auth(hevo.NewSourceAuth().ClientID("cid").ClientSecret("csecret").RefreshToken("rtoken"))).
		DestinationID(7).
		DestinationPrefix("hubspot_").
		ReplicationType("HISTORICAL_AND_INCREMENTAL").
		LoadMode("MERGE").
		SchemaEvolution("ALLOW_ALL").
		Schedule(hevo.NewSchedule().SyncType("SCHEDULED").FrequencyMinutes(60)).
		LatencyAlert(hevo.NewLatencyAlert().Enabled(true).ThresholdMinutes(120)).
		FailureHandlingPolicy(hevo.NewFailureHandlingPolicy().ObjectFailureLevel("PARTIAL").ObjectFailureThreshold(5)).
		CreateDisabled(true).
		IncludeAllObjects(false).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	assertPipeline(t, response.Data)
}

func TestCreatePipelineOmitsUnsetOptions(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodPost, "/api/v2/pipelines").ThenCall(func(req *http.Request) (*http.Response, error) {
		testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{"name": "p"})

		return mock.NewResponse(req, http.StatusOK, pipelineResponse("Success")), nil
	})

	if _, err := client.CreatePipeline().Name("p").Do(t.Context()); err != nil {
		t.Fatal(err)
	}
}
