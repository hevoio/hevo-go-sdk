//go:build e2e

package e2e_test

import (
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestPipelinesReadE2E(t *testing.T) {
	t.Parallel()

	client := testutils.E2EClient(t)

	list, err := client.ListPipelines().Limit(1).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(list.Data) == 0 {
		t.Skip("the account has no Pipelines to read")
	}

	id := list.Data[0].ID

	details, err := client.GetPipeline().PipelineID(id).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, details.Data.ID, id)

	if _, err := client.ListPipelineObjects().PipelineID(id).Limit(5).Do(t.Context()); err != nil {
		t.Fatal(err)
	}

	if _, err := client.ListPipelineJobs().PipelineID(id).Limit(5).Do(t.Context()); err != nil {
		t.Fatal(err)
	}
}
