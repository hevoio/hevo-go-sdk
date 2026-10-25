package metadata

import (
	"context"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListDestinationConnectorsService lists the connectors a Destination can use.
type ListDestinationConnectorsService struct {
	httputils.HTTPService

	limit  *int
	cursor *string
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListDestinationConnectorsService) Limit(value int) *ListDestinationConnectorsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListDestinationConnectorsService) Cursor(value string) *ListDestinationConnectorsService {
	s.cursor = &value

	return s
}

// Do reads one page.
func (s *ListDestinationConnectorsService) Do(ctx context.Context) (ConnectorsListResponse, error) {
	var response ConnectorsListResponse

	queries := httputils.PageQueries(s.limit, s.cursor)
	err := s.HTTPService.Do(ctx, http.MethodGet, destinationConnectorsPath, nil, queries, &response)

	return response, err
}
