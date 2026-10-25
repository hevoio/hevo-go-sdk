package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// SyncPipelineNowService starts an ingestion now instead of at the next scheduled time.
type SyncPipelineNowService struct {
	httputils.HTTPService

	pipelineID *int64
}

// PipelineID sets the Pipeline.
func (s *SyncPipelineNowService) PipelineID(value int64) *SyncPipelineNowService {
	s.pipelineID = &value

	return s
}

// Do sends the request.
func (s *SyncPipelineNowService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodPost, pipelinePath(*s.pipelineID, "/actions/sync-now"), nil, nil, &response)

	return response, err
}
