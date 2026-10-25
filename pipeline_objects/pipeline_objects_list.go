package pipelineobjects

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListPipelineObjectsService lists a Pipeline's objects one page at a time.
type ListPipelineObjectsService struct {
	httputils.HTTPService

	pipelineID *int64
	limit      *int
	cursor     *string
	statuses   []string
}

// PipelineID sets the Pipeline.
func (s *ListPipelineObjectsService) PipelineID(value int64) *ListPipelineObjectsService {
	s.pipelineID = &value

	return s
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListPipelineObjectsService) Limit(value int) *ListPipelineObjectsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListPipelineObjectsService) Cursor(value string) *ListPipelineObjectsService {
	s.cursor = &value

	return s
}

// Statuses keeps only objects in these statuses, for example "ACTIVE".
func (s *ListPipelineObjectsService) Statuses(values ...string) *ListPipelineObjectsService {
	s.statuses = values

	return s
}

// Do reads one page.
func (s *ListPipelineObjectsService) Do(ctx context.Context) (PipelineObjectsListResponse, error) {
	var response PipelineObjectsListResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	queries := httputils.PageQueries(s.limit, s.cursor)
	for _, status := range s.statuses {
		queries.Add("status", status)
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, objectsPath(*s.pipelineID, ""), nil, queries, &response)

	return response, err
}
