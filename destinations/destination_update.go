package destinations

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// UpdateDestinationService changes a Destination's name or settings. Only the fields set are
// sent; a Config holds only the keys to change.
type UpdateDestinationService struct {
	httputils.HTTPService

	destinationID *int64
	name          *string
	config        map[string]any
}

// DestinationID sets the Destination to update.
func (s *UpdateDestinationService) DestinationID(value int64) *UpdateDestinationService {
	s.destinationID = &value

	return s
}

// Name renames the Destination.
func (s *UpdateDestinationService) Name(value string) *UpdateDestinationService {
	s.name = &value

	return s
}

// Config sets the connector settings to change.
func (s *UpdateDestinationService) Config(value map[string]any) *UpdateDestinationService {
	s.config = value

	return s
}

// Do updates the Destination.
func (s *UpdateDestinationService) Do(ctx context.Context) (DestinationDetailsResponse, error) {
	var response DestinationDetailsResponse

	if s.destinationID == nil {
		return response, errors.New("hevo: missing required destinationID")
	}

	url := fmt.Sprintf("/api/v2/destinations/%d", *s.destinationID)
	req := &destinationUpdateRequest{Name: s.name, Config: s.config}
	err := s.HTTPService.Do(ctx, http.MethodPatch, url, req, nil, &response)

	return response, err
}
