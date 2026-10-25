package pipelines

import (
	"context"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// CreatePipelineService creates a Pipeline from a new Source into an existing Destination. The
// Pipeline is created disabled unless CreateDisabled(false) is set; enable it with the client's
// EnablePipeline once its objects are selected.
type CreatePipelineService struct {
	httputils.HTTPService

	name                  *string
	source                *Source
	destinationID         *int64
	destinationPrefix     *string
	replicationType       *string
	loadMode              *string
	schemaEvolution       *string
	schedule              *Schedule
	latencyAlert          *LatencyAlert
	failureHandlingPolicy *FailureHandlingPolicy
	createDisabled        *bool
	includeAllObjects     *bool
}

// Name sets the Pipeline's display name, unique within the team.
func (s *CreatePipelineService) Name(value string) *CreatePipelineService {
	s.name = &value

	return s
}

// Source sets the Source to create with the Pipeline.
func (s *CreatePipelineService) Source(value *Source) *CreatePipelineService {
	s.source = value

	return s
}

// DestinationID sets the Destination to load into.
func (s *CreatePipelineService) DestinationID(value int64) *CreatePipelineService {
	s.destinationID = &value

	return s
}

// DestinationPrefix sets the prefix of the tables created in the Destination.
func (s *CreatePipelineService) DestinationPrefix(value string) *CreatePipelineService {
	s.destinationPrefix = &value

	return s
}

// ReplicationType sets HISTORICAL_AND_INCREMENTAL or INCREMENTAL_ONLY.
func (s *CreatePipelineService) ReplicationType(value string) *CreatePipelineService {
	s.replicationType = &value

	return s
}

// LoadMode sets MERGE, APPEND or HISTORY_TRACKING.
func (s *CreatePipelineService) LoadMode(value string) *CreatePipelineService {
	s.loadMode = &value

	return s
}

// SchemaEvolution sets ALLOW_ALL or BLOCK_ALL: whether new objects and fields are replicated.
func (s *CreatePipelineService) SchemaEvolution(value string) *CreatePipelineService {
	s.schemaEvolution = &value

	return s
}

// Schedule sets when the Pipeline ingests.
func (s *CreatePipelineService) Schedule(value *Schedule) *CreatePipelineService {
	s.schedule = value

	return s
}

// LatencyAlert sets the latency alert.
func (s *CreatePipelineService) LatencyAlert(value *LatencyAlert) *CreatePipelineService {
	s.latencyAlert = value

	return s
}

// FailureHandlingPolicy sets how many objects may fail before a load fails.
func (s *CreatePipelineService) FailureHandlingPolicy(value *FailureHandlingPolicy) *CreatePipelineService {
	s.failureHandlingPolicy = value

	return s
}

// CreateDisabled leaves the Pipeline idle after creation. The API defaults it to true.
func (s *CreatePipelineService) CreateDisabled(value bool) *CreatePipelineService {
	s.createDisabled = &value

	return s
}

// IncludeAllObjects selects every object of the Source on creation.
func (s *CreatePipelineService) IncludeAllObjects(value bool) *CreatePipelineService {
	s.includeAllObjects = &value

	return s
}

// Do creates the Pipeline.
func (s *CreatePipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	err := s.HTTPService.Do(ctx, http.MethodPost, "/api/v2/pipelines", s.request(), nil, &response)

	return response, err
}

func (s *CreatePipelineService) request() *pipelineCreateRequest {
	req := &pipelineCreateRequest{
		Name:                  s.name,
		Source:                s.source.request(),
		DestinationID:         s.destinationID,
		DestinationPrefix:     s.destinationPrefix,
		ReplicationType:       s.replicationType,
		LoadMode:              s.loadMode,
		SchemaEvolution:       s.schemaEvolution,
		Schedule:              s.schedule.request(),
		LatencyAlert:          s.latencyAlert.request(),
		FailureHandlingPolicy: s.failureHandlingPolicy.request(),
	}

	if s.createDisabled != nil || s.includeAllObjects != nil {
		req.Options = &pipelineCreateOptionsRequest{
			CreateDisabled:    s.createDisabled,
			IncludeAllObjects: s.includeAllObjects,
		}
	}

	return req
}
