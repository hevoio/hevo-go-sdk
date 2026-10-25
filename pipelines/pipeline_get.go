package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetPipelineService reads one Pipeline.
type GetPipelineService struct {
	httputils.HTTPService

	pipelineID *int64
}

// PipelineID sets the Pipeline.
func (s *GetPipelineService) PipelineID(value int64) *GetPipelineService {
	s.pipelineID = &value

	return s
}

// Do sends the request.
func (s *GetPipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, pipelinePath(*s.pipelineID, ""), nil, nil, &response)

	return response, err
}
