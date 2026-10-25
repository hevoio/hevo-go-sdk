//go:build e2e

package e2e_test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hevoio/hevo-go"
	testutils "github.com/hevoio/hevo-go/test_utils"
)

// e2eDestination reads HEVO_E2E_DESTINATION: {"connector_id": "...", "config": {...}} with
// working credentials, because creating a Destination tests its connection.
func e2eDestination(t *testing.T) (string, map[string]any) {
	t.Helper()

	raw := os.Getenv("HEVO_E2E_DESTINATION")
	if raw == "" {
		t.Skip("set HEVO_E2E_DESTINATION to run the Destination lifecycle test")
	}

	var spec struct {
		ConnectorID string         `json:"connector_id"`
		Config      map[string]any `json:"config"`
	}
	if err := json.Unmarshal([]byte(raw), &spec); err != nil {
		t.Fatalf("HEVO_E2E_DESTINATION: %v", err)
	}

	return spec.ConnectorID, spec.Config
}

func TestDestinationLifecycleE2E(t *testing.T) {
	t.Parallel()

	client := testutils.E2EClient(t)
	connectorID, config := e2eDestination(t)
	name := "hevo-go-e2e-" + strconv.FormatInt(time.Now().UnixNano(), 36)

	created, err := client.CreateDestination().Name(name).ConnectorID(connectorID).Config(config).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	id := created.Data.ID

	t.Cleanup(func() {
		if _, err := client.DeleteDestination().DestinationID(id).Do(t.Context()); err != nil && !hevo.IsNotFound(err) {
			t.Errorf("deleting Destination %d: %v", id, err)
		}
	})

	testutils.AssertEqual(t, created.Data.Name, name)

	read, err := client.GetDestination().DestinationID(id).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, read.Data.ConnectorID, connectorID)

	updated, err := client.UpdateDestination().DestinationID(id).Name(name + "-renamed").Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, updated.Data.Name, name+"-renamed")
}
