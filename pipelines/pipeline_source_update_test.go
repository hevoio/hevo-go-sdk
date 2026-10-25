package pipelines_test

import (
	"net/http"
	"testing"

	"github.com/hevoio/hevo-go"
	testutils "github.com/hevoio/hevo-go/test_utils"
	"github.com/hevoio/hevo-go/tests/mock"
)

func TestUpdatePipelineSourceMock(t *testing.T) {
	t.Parallel()

	client, mockClient := testutils.CreateTestClient()
	handler := mockClient.When(http.MethodPatch, "/api/v2/pipelines/42/sources").ThenCall(
		func(req *http.Request) (*http.Response, error) {
			testutils.AssertEqual(t, testutils.RequestBodyToJSON(t, req), map[string]any{
				"config": map[string]any{"private_app_access_token": "pat-new"},
				"auth": map[string]any{
					"client_access": map[string]any{"client_id": "cid", "developer_token": "dev"},
					"access_token":  "atoken",
				},
			})

			return mock.NewResponse(req, http.StatusAccepted, pipelineResponse("Accepted")), nil
		})

	_, err := client.UpdatePipelineSource().
		PipelineID(42).
		Config(map[string]any{"private_app_access_token": "pat-new"}).
		Auth(hevo.NewSourceAuth().ClientID("cid").ClientAccessValue("developer_token", "dev").AccessToken("atoken")).
		Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, handler.Interactions, 1)
}
