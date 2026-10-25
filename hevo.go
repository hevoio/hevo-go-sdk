package hevo

import (
	httputils "github.com/hevoio/hevo-go/http_utils"
	pipelineobjects "github.com/hevoio/hevo-go/pipeline_objects"
	"github.com/hevoio/hevo-go/pipelines"
)

// NewSource starts the Source of a new Pipeline, for CreatePipeline.
func NewSource() *pipelines.Source {
	return &pipelines.Source{}
}

// NewSourceAuth starts the OAuth app credentials of a Source.
func NewSourceAuth() *pipelines.SourceAuth {
	return &pipelines.SourceAuth{}
}

// NewSchedule starts a Pipeline schedule.
func NewSchedule() *pipelines.Schedule {
	return &pipelines.Schedule{}
}

// NewLatencyAlert starts a Pipeline latency alert.
func NewLatencyAlert() *pipelines.LatencyAlert {
	return &pipelines.LatencyAlert{}
}

// NewFailureHandlingPolicy starts a Pipeline failure handling policy.
func NewFailureHandlingPolicy() *pipelines.FailureHandlingPolicy {
	return &pipelines.FailureHandlingPolicy{}
}

// NewObjectSelection starts one object's change, for UpdatePipelineObjects.
func NewObjectSelection() *pipelineobjects.ObjectSelection {
	return &pipelineobjects.ObjectSelection{}
}

// Debug turns on printing every request and response to stderr.
func Debug(enable bool) {
	httputils.Debug(enable)
}

// DebugAuth also prints the Authorization header while Debug is on.
func DebugAuth(enable bool) {
	httputils.DebugAuth(enable)
}

// IsNotFound reports whether err is a 404, for example a resource deleted outside the caller.
func IsNotFound(err error) bool { return httputils.IsNotFound(err) }

// IsConflict reports whether err is a 409: the Pipeline is mid-transition or has an edit open.
func IsConflict(err error) bool { return httputils.IsConflict(err) }

// IsRateLimited reports whether err is a 429 that outlasted the client's retries.
func IsRateLimited(err error) bool { return httputils.IsRateLimited(err) }
