package jobs

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListPipelineJobsService lists a Pipeline's jobs, newest first, one page at a time.
type ListPipelineJobsService struct {
	httputils.HTTPService

	pipelineID *int64
	limit      *int
	cursor     *string
	statuses   []string
	types      []string
}

// PipelineID sets the Pipeline.
func (s *ListPipelineJobsService) PipelineID(value int64) *ListPipelineJobsService {
	s.pipelineID = &value

	return s
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListPipelineJobsService) Limit(value int) *ListPipelineJobsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListPipelineJobsService) Cursor(value string) *ListPipelineJobsService {
	s.cursor = &value

	return s
}

// Statuses keeps only jobs in these statuses, for example "IN_PROGRESS".
func (s *ListPipelineJobsService) Statuses(values ...string) *ListPipelineJobsService {
	s.statuses = values

	return s
}

// Types keeps only jobs of these types, for example "HISTORICAL".
func (s *ListPipelineJobsService) Types(values ...string) *ListPipelineJobsService {
	s.types = values

	return s
}

// Do reads one page.
func (s *ListPipelineJobsService) Do(ctx context.Context) (JobsListResponse, error) {
	var response JobsListResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	queries := httputils.PageQueries(s.limit, s.cursor)
	for _, status := range s.statuses {
		queries.Add("status", status)
	}

	for _, jobType := range s.types {
		queries.Add("type", jobType)
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, jobsPath(*s.pipelineID), nil, queries, &response)

	return response, err
}
