package pipelineobjects_test

import (
	"net/http"
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestGetPipelineObjectMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/objects/public.orders").ThenReply(http.StatusOK,
		`{"code": "Success", "data": `+objectJSON+`}`)

	response, err := client.GetPipelineObject().PipelineID(42).ObjectName("public.orders").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, response.Data.Selected, true)
	testutils.AssertEqual(t, response.Data.Replication.Status, "REPLICATED")
}

func TestGetPipelineObjectEscapesTheName(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodGet, "/api/v2/pipelines/42/objects/a b/c").ThenReply(http.StatusOK,
		`{"code": "Success", "data": `+objectJSON+`}`)

	if _, err := client.GetPipelineObject().PipelineID(42).ObjectName("a b/c").Do(t.Context()); err != nil {
		t.Fatal(err)
	}

	interactions := mockClient.Interactions()
	testutils.AssertEqual(t, handler.Interactions, 1)
	testutils.AssertEqual(t, interactions[0].Request.URL.RawPath, "/api/v2/pipelines/42/objects/a%20b%2Fc")
}
