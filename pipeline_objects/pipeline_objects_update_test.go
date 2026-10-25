package pipelineobjects_test

import (
	"net/http"
	"testing"

	"github.com/hevoio/hevo-go"
	pipelineobjects "github.com/hevoio/hevo-go/pipeline_objects"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestUpdatePipelineObjectsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPatch, "/api/v2/pipelines/42/objects").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"selection_mode": "EXCLUDE_ALL",
				"objects": []any{
					map[string]any{"name": "public.orders", "selected": true, "load_mode": "APPEND"},
					map[string]any{"name": "public.users", "selected": true, "exclude_fields": []any{"ssn"}},
				},
			})

			return mock.NewResponse(req, http.StatusAccepted, acceptedPipeline), nil
		})

	response, err := client.UpdatePipelineObjects().
		PipelineID(42).
		SelectionMode(pipelineobjects.SelectionModeExcludeAll).
		Object(hevo.NewObjectSelection().Name("public.orders").Selected(true).LoadMode("APPEND")).
		Object(hevo.NewObjectSelection().Name("public.users").Selected(true).ExcludeFields("ssn")).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data.State, "Initializing")
}

func TestUpdatePipelineObjectsSendsDeselection(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	mockClient.When(http.MethodPatch, "/api/v2/pipelines/42/objects").ThenCall(func(req *http.Request) (*http.Response, error) {
		testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
			"objects": []any{map[string]any{"name": "public.orders", "selected": false}},
		})

		return mock.NewResponse(req, http.StatusAccepted, acceptedPipeline), nil
	})

	_, err := client.UpdatePipelineObjects().
		PipelineID(42).
		Object(hevo.NewObjectSelection().Name("public.orders").Selected(false)).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}
}
