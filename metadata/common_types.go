// Package metadata reads the connectors Hevo offers for Sources and Destinations, and the JSON
// Schema of each connector's config.
package metadata

import (
	"net/url"

	"github.com/hevoio/hevo-go/common"
)

// Connector is one source or destination connector.
type Connector struct {
	ConnectorID      string  `json:"connector_id"`
	Type             string  `json:"type"`
	DisplayName      string  `json:"display_name"`
	Description      *string `json:"description"`
	DocumentationURL *string `json:"documentation_url"`
	InBeta           bool    `json:"in_beta"`
}

// ConnectorDetails is a connector with the JSON Schema (draft 2020-12) of the body that creates
// it. Credentials are marked writeOnly in the schema.
type ConnectorDetails struct {
	ConnectorID string         `json:"connector_id"`
	Type        string         `json:"type"`
	DisplayName string         `json:"display_name"`
	Schema      map[string]any `json:"schema"`
}

// ConnectorsListResponse is one page of connectors.
type ConnectorsListResponse struct {
	common.CommonResponse
	common.ListResponse
	Data []Connector `json:"data"`
}

// ConnectorDetailsResponse is one connector with its schema.
type ConnectorDetailsResponse struct {
	common.CommonResponse
	Data ConnectorDetails `json:"data"`
}

const (
	sourceConnectorsPath      = "/api/v2/metadata/sources/connectors"
	destinationConnectorsPath = "/api/v2/metadata/destinations/connectors"
)

func connectorPath(base, connectorID string) string {
	return base + "/" + url.PathEscape(connectorID)
}
