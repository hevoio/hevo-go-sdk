package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// DeletePipelineService deletes a Pipeline. Deletion runs asynchronously: the response is Accepted and the Pipeline reports Deleting until it is gone.
type DeletePipelineService struct {
	httputils.HTTPService

	pipelineID *int64
}

// PipelineID sets the Pipeline.
func (s *DeletePipelineService) PipelineID(value int64) *DeletePipelineService {
	s.pipelineID = &value

	return s
}

// Do sends the request.
func (s *DeletePipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodDelete, pipelinePath(*s.pipelineID, ""), nil, nil, &response)

	return response, err
}
