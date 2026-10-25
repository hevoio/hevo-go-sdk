package pipelines_test

import (
	"testing"

	"github.com/hevoio/hevo-go/pipelines"
	testutils "github.com/hevoio/hevo-go/test_utils"
)

const pipelineJSON = `{
	"id": 42,
	"name": "HubSpot to Snowflake",
	"state": "Initialized",
	"disabled": false,
	"health": {"value": "HEALTHY", "display_name": "Healthy", "severity": "success"},
	"source": {"name": "HubSpot", "connector_id": "hubspot", "type": "HUBSPOT"},
	"destination_id": 7,
	"destination": {"name": "Snowflake", "connector_id": "snowflake", "type": "SNOWFLAKE", "status": "CREATED"},
	"destination_prefix": "hubspot_",
	"load_mode": "MERGE",
	"schema_evolution": "ALLOW_ALL",
	"replication_type": "HISTORICAL_AND_INCREMENTAL",
	"schedule": {"sync_type": "SCHEDULED", "frequency_minutes": 60},
	"latency_alert": {"enabled": true, "threshold_minutes": 120},
	"failure_handling_policy": {"object_failure_level": "PARTIAL", "object_failure_threshold": 5},
	"source_schema_status": {"status": "COMPLETED", "refreshed_at": 1760000000000},
	"runs": {
		"last_completed": {"job_id": "j-1", "type": "INCREMENTAL", "status": "COMPLETED", "created_at": 1760000000000},
		"active": [],
		"next_run": {"kind": "ON_SCHEDULE", "at": 1760003600000},
		"last_seven_days": {"incremental": {"ingested": 10, "loaded": 9, "failed": 1}}
	},
	"created_at": 1760000000000,
	"created_by": "ops@acme.com",
	"updated_at": 1760000001000,
	"updated_by": "ops@acme.com"
}`

func pipelineResponse(code string) string {
	return `{"code": "` + code + `", "data": ` + pipelineJSON + `}`
}

func assertPipeline(t *testing.T, p pipelines.PipelineDetails) {
	t.Helper()

	testutils.AssertEqual(t, p.ID, int64(42))
	testutils.AssertEqual(t, p.Name, "HubSpot to Snowflake")
	testutils.AssertEqual(t, p.State, "Initialized")
	testutils.AssertEqual(t, p.Health.Value, "HEALTHY")
	testutils.AssertEqual(t, p.Source, pipelines.PipelineSource{Name: "HubSpot", ConnectorID: "hubspot", Type: "HUBSPOT"})
	testutils.AssertEqual(t, p.DestinationID, int64(7))
	testutils.AssertEqual(t, p.Destination.Status, "CREATED")
	testutils.AssertEqual(t, p.DestinationPrefix, "hubspot_")
	testutils.AssertEqual(t, p.LoadMode, "MERGE")
	testutils.AssertEqual(t, p.SchemaEvolution, "ALLOW_ALL")
	testutils.AssertEqual(t, p.ReplicationType, "HISTORICAL_AND_INCREMENTAL")
	testutils.AssertEqual(t, p.Schedule, pipelines.ScheduleResponse{SyncType: "SCHEDULED", FrequencyMinutes: 60})
	testutils.AssertEqual(t, p.LatencyAlert, pipelines.LatencyAlertResponse{Enabled: true, ThresholdMinutes: 120})
	testutils.AssertEqual(t, p.FailureHandlingPolicy,
		pipelines.FailureHandlingPolicyResponse{ObjectFailureLevel: "PARTIAL", ObjectFailureThreshold: 5})
	testutils.AssertEqual(t, p.SourceSchemaStatus.Status, "COMPLETED")
	testutils.AssertEqual(t, p.Runs.LastCompleted.JobID, "j-1")
	testutils.AssertEqual(t, p.Runs.NextRun.Kind, "ON_SCHEDULE")
	testutils.AssertEqual(t, p.Runs.LastSevenDays.Incremental.Failed, int64(1))
	testutils.AssertEqual(t, p.UpdatedAt, int64(1760000001000))
}
