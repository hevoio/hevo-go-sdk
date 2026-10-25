package pipelines

import (
	"context"
	"net/http"
	"strconv"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// ListPipelinesService lists the team's Pipelines one page at a time, optionally filtered.
type ListPipelinesService struct {
	httputils.HTTPService

	limit         *int
	cursor        *string
	states        []string
	connectorID   *string
	destinationID *int64
}

// Limit sets the page size, from 1 to 500. The default is 100.
func (s *ListPipelinesService) Limit(value int) *ListPipelinesService {
	s.limit = &value

	return s
}

// Cursor asks for the page after the one whose NextCursor it is.
func (s *ListPipelinesService) Cursor(value string) *ListPipelinesService {
	s.cursor = &value

	return s
}

// States keeps only Pipelines in these states, for example "Initialized" or "Disabled".
func (s *ListPipelinesService) States(values ...string) *ListPipelinesService {
	s.states = values

	return s
}

// ConnectorID keeps only Pipelines whose Source uses this connector.
func (s *ListPipelinesService) ConnectorID(value string) *ListPipelinesService {
	s.connectorID = &value

	return s
}

// DestinationID keeps only Pipelines that load into this Destination.
func (s *ListPipelinesService) DestinationID(value int64) *ListPipelinesService {
	s.destinationID = &value

	return s
}

// Do reads one page.
func (s *ListPipelinesService) Do(ctx context.Context) (PipelinesListResponse, error) {
	var response PipelinesListResponse

	queries := httputils.PageQueries(s.limit, s.cursor)
	for _, state := range s.states {
		queries.Add("state", state)
	}

	if s.connectorID != nil {
		queries.Set("connector_id", *s.connectorID)
	}

	if s.destinationID != nil {
		queries.Set("destination_id", strconv.FormatInt(*s.destinationID, 10))
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, "/api/v2/pipelines", nil, queries, &response)

	return response, err
}
