package jobs

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListPipelineJobObjectsService lists what one job did for each object, one page at a time.
type ListPipelineJobObjectsService struct {
	httputils.HTTPService

	pipelineID *int64
	jobID      *string
	limit      *int
	cursor     *string
}

// PipelineID sets the Pipeline.
func (s *ListPipelineJobObjectsService) PipelineID(value int64) *ListPipelineJobObjectsService {
	s.pipelineID = &value

	return s
}

// JobID sets the job.
func (s *ListPipelineJobObjectsService) JobID(value string) *ListPipelineJobObjectsService {
	s.jobID = &value

	return s
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListPipelineJobObjectsService) Limit(value int) *ListPipelineJobObjectsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListPipelineJobObjectsService) Cursor(value string) *ListPipelineJobObjectsService {
	s.cursor = &value

	return s
}

// Do sends the request.
func (s *ListPipelineJobObjectsService) Do(ctx context.Context) (JobObjectsListResponse, error) {
	var response JobObjectsListResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	if s.jobID == nil {
		return response, errors.New("hevo: missing required jobID")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, jobPath(*s.pipelineID, *s.jobID, "/objects"), nil, httputils.PageQueries(s.limit, s.cursor), &response)

	return response, err
}
