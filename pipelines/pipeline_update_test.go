package pipelines_test

import (
	"net/http"
	"testing"

	"github.com/hevoio/hevo-go"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestUpdatePipelineMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPatch, "/api/v2/pipelines/42").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"name":             "Renamed",
				"load_mode":        "APPEND",
				"schema_evolution": "BLOCK_ALL",
				"schedule":         map[string]any{"sync_type": "ON_DEMAND"},
			})

			return mock.NewResponse(req, http.StatusAccepted, pipelineResponse("Accepted")), nil
		})

	response, err := client.UpdatePipeline().
		PipelineID(42).
		Name("Renamed").
		LoadMode("APPEND").
		SchemaEvolution("BLOCK_ALL").
		Schedule(hevo.NewSchedule().SyncType("ON_DEMAND")).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Code, "Accepted")
}

func TestUpdatePipelineConflict(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodPatch, "/api/v2/pipelines/42").ThenReply(http.StatusConflict,
		`{"code": "RESOURCE_CONFLICT", "message": "The Pipeline has an edit in progress."}`)

	_, err := client.UpdatePipeline().PipelineID(42).Name("x").Do(t.Context())
	if !hevo.IsConflict(err) {
		t.Fatalf("err = %v, want a 409", err)
	}
}
