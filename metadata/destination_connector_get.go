package metadata

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetDestinationConnectorService reads a destination connector with the JSON Schema of its config.
type GetDestinationConnectorService struct {
	httputils.HTTPService

	connectorID *string
}

// ConnectorID sets the connector, for example "snowflake".
func (s *GetDestinationConnectorService) ConnectorID(value string) *GetDestinationConnectorService {
	s.connectorID = &value

	return s
}

// Do reads the connector.
func (s *GetDestinationConnectorService) Do(ctx context.Context) (ConnectorDetailsResponse, error) {
	var response ConnectorDetailsResponse

	if s.connectorID == nil {
		return response, errors.New("hevo: missing required connectorID")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, connectorPath(destinationConnectorsPath, *s.connectorID), nil, nil, &response)

	return response, err
}
