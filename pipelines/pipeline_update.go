package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// UpdatePipelineService changes a Pipeline's settings. Only the fields set are sent. The edit runs
// asynchronously: a replicating Pipeline reports Initializing until it is applied. A 409 means
// the Pipeline is mid-transition or has an edit open in the Hevo UI; retry once it settles.
type UpdatePipelineService struct {
	httputils.HTTPService

	pipelineID            *int64
	name                  *string
	destinationPrefix     *string
	loadMode              *string
	schemaEvolution       *string
	schedule              *Schedule
	latencyAlert          *LatencyAlert
	failureHandlingPolicy *FailureHandlingPolicy
}

// PipelineID sets the Pipeline to update.
func (s *UpdatePipelineService) PipelineID(value int64) *UpdatePipelineService {
	s.pipelineID = &value

	return s
}

// Name renames the Pipeline.
func (s *UpdatePipelineService) Name(value string) *UpdatePipelineService {
	s.name = &value

	return s
}

// DestinationPrefix changes the prefix of the tables created in the Destination.
func (s *UpdatePipelineService) DestinationPrefix(value string) *UpdatePipelineService {
	s.destinationPrefix = &value

	return s
}

// LoadMode sets MERGE, APPEND or HISTORY_TRACKING.
func (s *UpdatePipelineService) LoadMode(value string) *UpdatePipelineService {
	s.loadMode = &value

	return s
}

// SchemaEvolution sets ALLOW_ALL or BLOCK_ALL.
func (s *UpdatePipelineService) SchemaEvolution(value string) *UpdatePipelineService {
	s.schemaEvolution = &value

	return s
}

// Schedule changes when the Pipeline ingests.
func (s *UpdatePipelineService) Schedule(value *Schedule) *UpdatePipelineService {
	s.schedule = value

	return s
}

// LatencyAlert changes the latency alert.
func (s *UpdatePipelineService) LatencyAlert(value *LatencyAlert) *UpdatePipelineService {
	s.latencyAlert = value

	return s
}

// FailureHandlingPolicy changes how many objects may fail before a load fails.
func (s *UpdatePipelineService) FailureHandlingPolicy(value *FailureHandlingPolicy) *UpdatePipelineService {
	s.failureHandlingPolicy = value

	return s
}

// Do updates the Pipeline.
func (s *UpdatePipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	req := &pipelineUpdateRequest{
		Name:                  s.name,
		DestinationPrefix:     s.destinationPrefix,
		LoadMode:              s.loadMode,
		SchemaEvolution:       s.schemaEvolution,
		Schedule:              s.schedule.request(),
		LatencyAlert:          s.latencyAlert.request(),
		FailureHandlingPolicy: s.failureHandlingPolicy.request(),
	}
	err := s.HTTPService.Do(ctx, http.MethodPatch, pipelinePath(*s.pipelineID, ""), req, nil, &response)

	return response, err
}
