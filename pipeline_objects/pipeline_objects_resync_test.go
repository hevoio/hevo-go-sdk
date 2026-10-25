package pipelineobjects_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestResyncPipelineObjectsMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPost, "/api/v2/pipelines/42/objects/actions/resync").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"object_names":  []any{"public.orders"},
				"drop_and_load": true,
			})

			return mock.NewResponse(req, http.StatusAccepted, acceptedPipeline), nil
		})

	_, err := client.ResyncPipelineObjects().PipelineID(42).ObjectNames("public.orders").DropAndLoad(true).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
}

func TestResyncPipelineObjectsRequiresNames(t *testing.T) {
	t.Parallel()

	client, _ := testutils.CreateTestClient()
	if _, err := client.ResyncPipelineObjects().PipelineID(42).Do(t.Context()); err == nil {
		t.Fatal("want an error without object names")
	}
}
