//go:build e2e

package e2e_test

import (
	"testing"

	testutils "github.com/hevoio/hevo-go/test_utils"
)

func TestSourceConnectorsE2E(t *testing.T) {
	t.Parallel()

	client := testutils.E2EClient(t)

	list, err := client.ListSourceConnectors().Limit(5).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(list.Data) == 0 {
		t.Fatal("no source connectors listed")
	}

	details, err := client.GetSourceConnector().ConnectorID(list.Data[0].ConnectorID).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, details.Data.ConnectorID, list.Data[0].ConnectorID)

	if len(details.Data.Schema) == 0 {
		t.Error("connector schema is empty")
	}
}

func TestDestinationConnectorsE2E(t *testing.T) {
	t.Parallel()

	client := testutils.E2EClient(t)

	list, err := client.ListDestinationConnectors().Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(list.Data) == 0 {
		t.Fatal("no destination connectors listed")
	}

	if _, err := client.GetDestinationConnector().ConnectorID(list.Data[0].ConnectorID).Do(t.Context()); err != nil {
		t.Fatal(err)
	}
}
