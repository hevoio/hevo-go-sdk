package destinations

import (
	"context"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// CreateDestinationService creates a Destination. Config follows the connector's schema, from
// the client's GetDestinationConnector.
type CreateDestinationService struct {
	httputils.HTTPService

	name        *string
	connectorID *string
	destType    *string
	config      map[string]any
}

// Name sets the Destination's display name, unique within the team.
func (s *CreateDestinationService) Name(value string) *CreateDestinationService {
	s.name = &value

	return s
}

// ConnectorID sets the destination connector, for example "snowflake".
func (s *CreateDestinationService) ConnectorID(value string) *CreateDestinationService {
	s.connectorID = &value

	return s
}

// Type sets the destination type, for example "SNOWFLAKE". It is optional when ConnectorID
// identifies the connector.
func (s *CreateDestinationService) Type(value string) *CreateDestinationService {
	s.destType = &value

	return s
}

// Config sets the connector's settings, keyed as its schema names them.
func (s *CreateDestinationService) Config(value map[string]any) *CreateDestinationService {
	s.config = value

	return s
}

// Do creates the Destination.
func (s *CreateDestinationService) Do(ctx context.Context) (DestinationDetailsResponse, error) {
	var response DestinationDetailsResponse

	err := s.HTTPService.Do(ctx, http.MethodPost, "/api/v2/destinations", s.request(), nil, &response)

	return response, err
}

func (s *CreateDestinationService) request() *destinationCreateRequest {
	return &destinationCreateRequest{
		Name:        s.name,
		ConnectorID: s.connectorID,
		Type:        s.destType,
		Config:      configPointer(s.config),
	}
}

// configPointer keeps an empty config in the body, as {}; a config never set is left out.
func configPointer(config map[string]any) *map[string]any {
	if config == nil {
		return nil
	}

	return &config
}
