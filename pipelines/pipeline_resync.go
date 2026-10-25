package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ResyncPipelineService reloads every selected object from the start. It runs asynchronously.
type ResyncPipelineService struct {
	httputils.HTTPService

	pipelineID  *int64
	dropAndLoad *bool
}

// PipelineID sets the Pipeline.
func (s *ResyncPipelineService) PipelineID(value int64) *ResyncPipelineService {
	s.pipelineID = &value

	return s
}

// DropAndLoad drops the Destination tables before reloading them.
func (s *ResyncPipelineService) DropAndLoad(value bool) *ResyncPipelineService {
	s.dropAndLoad = &value

	return s
}

// Do sends the request.
func (s *ResyncPipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodPost, pipelinePath(*s.pipelineID, "/actions/resync"), &pipelineResyncRequest{DropAndLoad: s.dropAndLoad}, nil, &response)

	return response, err
}
