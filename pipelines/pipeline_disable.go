package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// DisablePipelineService pauses a Pipeline. It runs asynchronously: the Pipeline reports Disabling, then Disabled.
type DisablePipelineService struct {
	httputils.HTTPService

	pipelineID       *int64
	cancelActiveJobs *bool
}

// PipelineID sets the Pipeline.
func (s *DisablePipelineService) PipelineID(value int64) *DisablePipelineService {
	s.pipelineID = &value

	return s
}

// CancelActiveJobs cancels the jobs in progress instead of letting them finish.
func (s *DisablePipelineService) CancelActiveJobs(value bool) *DisablePipelineService {
	s.cancelActiveJobs = &value

	return s
}

// Do sends the request.
func (s *DisablePipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodPost, pipelinePath(*s.pipelineID, "/actions/disable"), &pipelineDisableRequest{CancelActiveJobs: s.cancelActiveJobs}, nil, &response)

	return response, err
}
