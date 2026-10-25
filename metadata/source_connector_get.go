package metadata

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetSourceConnectorService reads a source connector with the JSON Schema of its config.
type GetSourceConnectorService struct {
	httputils.HTTPService

	connectorID *string
}

// ConnectorID sets the connector, for example "hubspot".
func (s *GetSourceConnectorService) ConnectorID(value string) *GetSourceConnectorService {
	s.connectorID = &value

	return s
}

// Do reads the connector.
func (s *GetSourceConnectorService) Do(ctx context.Context) (ConnectorDetailsResponse, error) {
	var response ConnectorDetailsResponse

	if s.connectorID == nil {
		return response, errors.New("hevo: missing required connectorID")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, connectorPath(sourceConnectorsPath, *s.connectorID), nil, nil, &response)

	return response, err
}
