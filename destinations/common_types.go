// Package destinations manages Destinations: where Pipelines load data.
package destinations

import "github.com/hevoio/hevo-go/common"

// DestinationDetails is one Destination. Credentials in Config come back masked.
type DestinationDetails struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	ConnectorID string         `json:"connector_id"`
	Type        string         `json:"type"`
	Status      string         `json:"status"`
	Config      map[string]any `json:"config"`
	CreatedAt   int64          `json:"created_at"`
	CreatedBy   string         `json:"created_by"`
	UpdatedAt   int64          `json:"updated_at"`
	UpdatedBy   string         `json:"updated_by"`
}

// DestinationDetailsResponse is the response of create, details, update and delete.
type DestinationDetailsResponse struct {
	common.CommonResponse
	Data DestinationDetails `json:"data"`
}

// DestinationsListResponse is one page of Destinations.
type DestinationsListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []DestinationDetails `json:"data"`
}

type destinationCreateRequest struct {
	Name        *string         `json:"name,omitempty"`
	ConnectorID *string         `json:"connector_id,omitempty"`
	Type        *string         `json:"type,omitempty"`
	Config      *map[string]any `json:"config,omitempty"`
}

type destinationUpdateRequest struct {
	Name   *string        `json:"name,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}
