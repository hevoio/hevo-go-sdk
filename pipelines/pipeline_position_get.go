package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetPipelinePositionService reads how far replication has read each object, one page at a time.
type GetPipelinePositionService struct {
	httputils.HTTPService

	pipelineID *int64
	limit      *int
	cursor     *string
}

// PipelineID sets the Pipeline.
func (s *GetPipelinePositionService) PipelineID(value int64) *GetPipelinePositionService {
	s.pipelineID = &value

	return s
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *GetPipelinePositionService) Limit(value int) *GetPipelinePositionService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *GetPipelinePositionService) Cursor(value string) *GetPipelinePositionService {
	s.cursor = &value

	return s
}

// Do reads one page.
func (s *GetPipelinePositionService) Do(ctx context.Context) (PipelinePositionResponse, error) {
	var response PipelinePositionResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	queries := httputils.PageQueries(s.limit, s.cursor)
	err := s.HTTPService.Do(ctx, http.MethodGet, pipelinePath(*s.pipelineID, "/position"), nil, queries, &response)

	return response, err
}
