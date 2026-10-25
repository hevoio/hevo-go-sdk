// Package jobs reads the jobs of a Pipeline: each ingestion, historical load or resync it runs.
package jobs

import (
	"fmt"
	"net/url"

	"github.com/hevoio/hevo-go/common"
)

// Job is one run of a Pipeline.
type Job struct {
	ID             string     `json:"id"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	CreatedAt      int64      `json:"created_at"`
	UpdatedAt      int64      `json:"updated_at"`
	EventsIngested int64      `json:"events_ingested"`
	EventsLoaded   int64      `json:"events_loaded"`
	EventsFailed   int64      `json:"events_failed"`
	ObjectsSuccess int64      `json:"objects_success"`
	ObjectsQueued  int64      `json:"objects_queued"`
	ObjectsSkipped int64      `json:"objects_skipped"`
	ObjectsFailed  int64      `json:"objects_failed"`
	Duration       int64      `json:"duration"`
	Latency        JobLatency `json:"latency"`
}

// JobLatency is the job's latency in milliseconds.
type JobLatency struct {
	Min  int64 `json:"min"`
	Max  int64 `json:"max"`
	Mean int64 `json:"mean"`
}

// JobObject is what one job did for one object.
type JobObject struct {
	ID             string  `json:"id"`
	ObjectName     string  `json:"object_name"`
	Status         string  `json:"status"`
	EventsIngested *int64  `json:"events_ingested"`
	EventsLoaded   *int64  `json:"events_loaded"`
	EventsFailed   *int64  `json:"events_failed"`
	Offset         *string `json:"offset"`
}

// JobResponse is one job.
type JobResponse struct {
	common.CommonResponse
	Data Job `json:"data"`
}

// JobsListResponse is one page of jobs.
type JobsListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []Job `json:"data"`
}

// JobObjectsListResponse is one page of a job's objects.
type JobObjectsListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []JobObject `json:"data"`
}

func jobsPath(pipelineID int64) string {
	return fmt.Sprintf("/api/v2/pipelines/%d/jobs", pipelineID)
}

func jobPath(pipelineID int64, jobID, suffix string) string {
	return jobsPath(pipelineID) + "/" + url.PathEscape(jobID) + suffix
}
