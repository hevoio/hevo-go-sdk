package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// UpdatePipelineSourceService changes the Pipeline's Source: its name, settings or OAuth
// credentials. Only the connector's editable fields may be sent; Config holds only the keys to
// change. It runs asynchronously, like UpdatePipelineService.
type UpdatePipelineSourceService struct {
	httputils.HTTPService

	pipelineID *int64
	name       *string
	config     map[string]any
	auth       *SourceAuth
}

// PipelineID sets the Pipeline whose Source to update.
func (s *UpdatePipelineSourceService) PipelineID(value int64) *UpdatePipelineSourceService {
	s.pipelineID = &value

	return s
}

// Name renames the Source.
func (s *UpdatePipelineSourceService) Name(value string) *UpdatePipelineSourceService {
	s.name = &value

	return s
}

// Config sets the connector settings to change.
func (s *UpdatePipelineSourceService) Config(value map[string]any) *UpdatePipelineSourceService {
	s.config = value

	return s
}

// Auth replaces the OAuth app credentials.
func (s *UpdatePipelineSourceService) Auth(value *SourceAuth) *UpdatePipelineSourceService {
	s.auth = value

	return s
}

// Do updates the Source.
func (s *UpdatePipelineSourceService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	req := &pipelineSourceUpdateRequest{Name: s.name, Config: s.config, Auth: s.auth.request()}
	err := s.HTTPService.Do(ctx, http.MethodPatch, pipelinePath(*s.pipelineID, "/sources"), req, nil, &response)

	return response, err
}
