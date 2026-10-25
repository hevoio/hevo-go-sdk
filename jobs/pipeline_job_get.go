package jobs

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetPipelineJobService reads one job.
type GetPipelineJobService struct {
	httputils.HTTPService

	pipelineID *int64
	jobID      *string
}

// PipelineID sets the Pipeline.
func (s *GetPipelineJobService) PipelineID(value int64) *GetPipelineJobService {
	s.pipelineID = &value

	return s
}

// JobID sets the job.
func (s *GetPipelineJobService) JobID(value string) *GetPipelineJobService {
	s.jobID = &value

	return s
}

// Do sends the request.
func (s *GetPipelineJobService) Do(ctx context.Context) (JobResponse, error) {
	var response JobResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	if s.jobID == nil {
		return response, errors.New("hevo: missing required jobID")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, jobPath(*s.pipelineID, *s.jobID, ""), nil, nil, &response)

	return response, err
}
