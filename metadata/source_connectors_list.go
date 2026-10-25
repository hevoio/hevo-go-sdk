package metadata

import (
	"context"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListSourceConnectorsService lists the connectors a Pipeline's Source can use.
type ListSourceConnectorsService struct {
	httputils.HTTPService

	limit      *int
	cursor     *string
	sourceType *string
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListSourceConnectorsService) Limit(value int) *ListSourceConnectorsService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListSourceConnectorsService) Cursor(value string) *ListSourceConnectorsService {
	s.cursor = &value

	return s
}

// Type keeps only the connector of this source type, for example "HUBSPOT".
func (s *ListSourceConnectorsService) Type(value string) *ListSourceConnectorsService {
	s.sourceType = &value

	return s
}

// Do reads one page.
func (s *ListSourceConnectorsService) Do(ctx context.Context) (ConnectorsListResponse, error) {
	var response ConnectorsListResponse

	queries := httputils.PageQueries(s.limit, s.cursor)
	if s.sourceType != nil {
		queries.Set("type", *s.sourceType)
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, sourceConnectorsPath, nil, queries, &response)

	return response, err
}
