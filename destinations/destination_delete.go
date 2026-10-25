package destinations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// DeleteDestinationService deletes a Destination. Deletion runs asynchronously; the response is
// Accepted and the Destination reports DELETING until it is gone.
type DeleteDestinationService struct {
	httputils.HTTPService

	destinationID   *int64
	confirmOverride *bool
}

// DestinationID sets the Destination to delete.
func (s *DeleteDestinationService) DestinationID(value int64) *DeleteDestinationService {
	s.destinationID = &value

	return s
}

// ConfirmOverride deletes the Destination even when Pipelines still load into it.
func (s *DeleteDestinationService) ConfirmOverride(value bool) *DeleteDestinationService {
	s.confirmOverride = &value

	return s
}

// Do deletes the Destination.
func (s *DeleteDestinationService) Do(ctx context.Context) (DestinationDetailsResponse, error) {
	var response DestinationDetailsResponse

	if s.destinationID == nil {
		return response, errors.New("hevo: missing required destinationID")
	}

	var queries url.Values
	if s.confirmOverride != nil {
		queries = url.Values{"confirm_override": {strconv.FormatBool(*s.confirmOverride)}}
	}

	path := fmt.Sprintf("/api/v2/destinations/%d", *s.destinationID)
	err := s.HTTPService.Do(ctx, http.MethodDelete, path, nil, queries, &response)

	return response, err
}
