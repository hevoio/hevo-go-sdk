package destinations

import (
	"context"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListDestinationsService lists the team's Destinations one page at a time, ordered by id.
type ListDestinationsService struct {
	httputils.HTTPService

	limit  *int
	cursor *string
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListDestinationsService) Limit(value int) *ListDestinationsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListDestinationsService) Cursor(value string) *ListDestinationsService {
	s.cursor = &value

	return s
}

// Do reads one page.
func (s *ListDestinationsService) Do(ctx context.Context) (DestinationsListResponse, error) {
	var response DestinationsListResponse

	queries := httputils.PageQueries(s.limit, s.cursor)
	err := s.HTTPService.Do(ctx, http.MethodGet, "/api/v2/destinations", nil, queries, &response)

	return response, err
}
