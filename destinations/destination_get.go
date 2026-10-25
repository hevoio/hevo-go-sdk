package destinations

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetDestinationService reads one Destination.
type GetDestinationService struct {
	httputils.HTTPService

	destinationID *int64
}

// DestinationID sets the Destination to read.
func (s *GetDestinationService) DestinationID(value int64) *GetDestinationService {
	s.destinationID = &value

	return s
}

// Do reads the Destination.
func (s *GetDestinationService) Do(ctx context.Context) (DestinationDetailsResponse, error) {
	var response DestinationDetailsResponse

	if s.destinationID == nil {
		return response, errors.New("hevo: missing required destinationID")
	}

	url := fmt.Sprintf("/api/v2/destinations/%d", *s.destinationID)
	err := s.HTTPService.Do(ctx, http.MethodGet, url, nil, nil, &response)

	return response, err
}
