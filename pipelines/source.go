package pipelines

// Source is the Source a new Pipeline reads from. Build it with the client's NewSource.
type Source struct {
	name        *string
	connectorID *string
	sourceType  *string
	config      map[string]any
	auth        *SourceAuth
}

type sourceRequest struct {
	Name        *string            `json:"name,omitempty"`
	ConnectorID *string            `json:"connector_id,omitempty"`
	Type        *string            `json:"type,omitempty"`
	Config      map[string]any     `json:"config,omitempty"`
	Auth        *sourceAuthRequest `json:"auth,omitempty"`
}

// Name sets the Source's display name, unique within the team.
func (s *Source) Name(value string) *Source {
	s.name = &value

	return s
}

// ConnectorID sets the source connector, for example "hubspot".
func (s *Source) ConnectorID(value string) *Source {
	s.connectorID = &value

	return s
}

// Type sets the source type, for example "HUBSPOT". It is optional when ConnectorID identifies
// the connector.
func (s *Source) Type(value string) *Source {
	s.sourceType = &value

	return s
}

// Config sets the connector's settings, keyed as its schema names them. Credentials are
// write-only: they are masked in every read.
func (s *Source) Config(value map[string]any) *Source {
	s.config = value

	return s
}

// Auth sets the OAuth app credentials for a Source that authenticates with your own OAuth app.
func (s *Source) Auth(value *SourceAuth) *Source {
	s.auth = value

	return s
}

// request returns the Source as the API expects it.
func (s *Source) request() *sourceRequest {
	if s == nil {
		return nil
	}

	return &sourceRequest{
		Name:        s.name,
		ConnectorID: s.connectorID,
		Type:        s.sourceType,
		Config:      s.config,
		Auth:        s.auth.request(),
	}
}
