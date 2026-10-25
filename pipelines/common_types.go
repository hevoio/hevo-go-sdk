// Package pipelines manages Pipelines: each one replicates a Source into a Destination.
package pipelines

import (
	"fmt"

	"github.com/hevoio/hevo-go/common"
)

// PipelineDetails is one Pipeline as the API returns it.
type PipelineDetails struct {
	ID                    int64                         `json:"id"`
	Name                  string                        `json:"name"`
	State                 string                        `json:"state"`
	Disabled              bool                          `json:"disabled"`
	Health                PipelineHealth                `json:"health"`
	Source                PipelineSource                `json:"source"`
	DestinationID         int64                         `json:"destination_id"`
	Destination           PipelineDestination           `json:"destination"`
	DestinationPrefix     string                        `json:"destination_prefix"`
	LoadMode              string                        `json:"load_mode"`
	SchemaEvolution       string                        `json:"schema_evolution"`
	ReplicationType       string                        `json:"replication_type"`
	Schedule              ScheduleResponse              `json:"schedule"`
	LatencyAlert          LatencyAlertResponse          `json:"latency_alert"`
	FailureHandlingPolicy FailureHandlingPolicyResponse `json:"failure_handling_policy"`
	SourceSchemaStatus    SourceSchemaStatus            `json:"source_schema_status"`
	Runs                  PipelineRuns                  `json:"runs"`
	CreatedAt             int64                         `json:"created_at"`
	CreatedBy             string                        `json:"created_by"`
	UpdatedAt             int64                         `json:"updated_at"`
	UpdatedBy             string                        `json:"updated_by"`
}

// PipelineHealth is the Pipeline's health as the Hevo UI shows it.
type PipelineHealth struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
	Severity    string `json:"severity"`
}

// PipelineSource identifies the Pipeline's Source.
type PipelineSource struct {
	Name        string `json:"name"`
	ConnectorID string `json:"connector_id"`
	Type        string `json:"type"`
}

// PipelineDestination identifies the Destination the Pipeline loads into.
type PipelineDestination struct {
	Name        string `json:"name"`
	ConnectorID string `json:"connector_id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
}

// SourceSchemaStatus is the state of the Source's object discovery.
type SourceSchemaStatus struct {
	Status      string `json:"status"`
	RefreshedAt int64  `json:"refreshed_at"`
}

// PipelineRuns summarises the Pipeline's jobs.
type PipelineRuns struct {
	LastCompleted *PipelineRun          `json:"last_completed"`
	Active        []PipelineRun         `json:"active"`
	NextRun       *PipelineNextRun      `json:"next_run"`
	LastSevenDays PipelineEventCountSet `json:"last_seven_days"`
}

// PipelineRun is one job of the Pipeline.
type PipelineRun struct {
	JobID     string `json:"job_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

// PipelineNextRun says when the next job starts.
type PipelineNextRun struct {
	Kind string `json:"kind"`
	At   int64  `json:"at"`
}

// PipelineEventCountSet counts the last seven days' events by kind of job.
type PipelineEventCountSet struct {
	Historical  PipelineEventCounts `json:"historical"`
	Incremental PipelineEventCounts `json:"incremental"`
	Refresher   PipelineEventCounts `json:"refresher"`
}

// PipelineEventCounts counts events ingested, loaded and failed.
type PipelineEventCounts struct {
	Ingested int64 `json:"ingested"`
	Loaded   int64 `json:"loaded"`
	Failed   int64 `json:"failed"`
}

// PipelineDetailsResponse is the response of every single-Pipeline call. Code is "Accepted"
// when the work runs asynchronously; Data then shows the Pipeline in a transitional state.
type PipelineDetailsResponse struct {
	common.CommonResponse
	Data PipelineDetails `json:"data"`
}

// PipelinesListResponse is one page of Pipelines.
type PipelinesListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []PipelineDetails `json:"data"`
}

// PipelinePosition is how far replication has read one object.
type PipelinePosition struct {
	FullyQualifiedName string  `json:"fully_qualified_name"`
	Name               string  `json:"name"`
	Position           *string `json:"position"`
}

// PipelinePositionResponse is one page of object positions.
type PipelinePositionResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []PipelinePosition `json:"data"`
}

type pipelineCreateRequest struct {
	Name                  *string                       `json:"name,omitempty"`
	Source                *sourceRequest                `json:"source,omitempty"`
	DestinationID         *int64                        `json:"destination_id,omitempty"`
	DestinationPrefix     *string                       `json:"destination_prefix,omitempty"`
	ReplicationType       *string                       `json:"replication_type,omitempty"`
	LoadMode              *string                       `json:"load_mode,omitempty"`
	SchemaEvolution       *string                       `json:"schema_evolution,omitempty"`
	Schedule              *scheduleRequest              `json:"schedule,omitempty"`
	LatencyAlert          *latencyAlertRequest          `json:"latency_alert,omitempty"`
	FailureHandlingPolicy *failureHandlingPolicyRequest `json:"failure_handling_policy,omitempty"`
	Options               *pipelineCreateOptionsRequest `json:"options,omitempty"`
}

type pipelineCreateOptionsRequest struct {
	CreateDisabled    *bool `json:"create_disabled,omitempty"`
	IncludeAllObjects *bool `json:"include_all_objects,omitempty"`
}

type pipelineUpdateRequest struct {
	Name                  *string                       `json:"name,omitempty"`
	DestinationPrefix     *string                       `json:"destination_prefix,omitempty"`
	LoadMode              *string                       `json:"load_mode,omitempty"`
	SchemaEvolution       *string                       `json:"schema_evolution,omitempty"`
	Schedule              *scheduleRequest              `json:"schedule,omitempty"`
	LatencyAlert          *latencyAlertRequest          `json:"latency_alert,omitempty"`
	FailureHandlingPolicy *failureHandlingPolicyRequest `json:"failure_handling_policy,omitempty"`
}

type pipelineSourceUpdateRequest struct {
	Name   *string            `json:"name,omitempty"`
	Config map[string]any     `json:"config,omitempty"`
	Auth   *sourceAuthRequest `json:"auth,omitempty"`
}

type pipelineDisableRequest struct {
	CancelActiveJobs *bool `json:"cancel_active_jobs,omitempty"`
}

type pipelineResyncRequest struct {
	DropAndLoad *bool `json:"drop_and_load,omitempty"`
}

// pipelinePath is the path of a Pipeline, or of one of its sub-resources when suffix is set.
func pipelinePath(pipelineID int64, suffix string) string {
	return fmt.Sprintf("/api/v2/pipelines/%d%s", pipelineID, suffix)
}
