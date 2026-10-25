package jobs

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// CancelPipelineJobService cancels a job in progress. The job reports CANCELING until it stops.
type CancelPipelineJobService struct {
	httputils.HTTPService

	pipelineID *int64
	jobID      *string
}

// PipelineID sets the Pipeline.
func (s *CancelPipelineJobService) PipelineID(value int64) *CancelPipelineJobService {
	s.pipelineID = &value

	return s
}

// JobID sets the job.
func (s *CancelPipelineJobService) JobID(value string) *CancelPipelineJobService {
	s.jobID = &value

	return s
}

// Do sends the request.
func (s *CancelPipelineJobService) Do(ctx context.Context) (JobResponse, error) {
	var response JobResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	if s.jobID == nil {
		return response, errors.New("hevo: missing required jobID")
	}

	err := s.HTTPService.Do(ctx, http.MethodPost, jobPath(*s.pipelineID, *s.jobID, "/actions/cancel"), nil, nil, &response)

	return response, err
}
