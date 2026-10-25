package pipelineobjects

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// GetPipelineObjectService reads one object of a Pipeline, with its fields.
type GetPipelineObjectService struct {
	httputils.HTTPService

	pipelineID *int64
	objectName *string
}

// PipelineID sets the Pipeline.
func (s *GetPipelineObjectService) PipelineID(value int64) *GetPipelineObjectService {
	s.pipelineID = &value

	return s
}

// ObjectName sets the object, as its name appears in the list of objects.
func (s *GetPipelineObjectService) ObjectName(value string) *GetPipelineObjectService {
	s.objectName = &value

	return s
}

// Do reads the object.
func (s *GetPipelineObjectService) Do(ctx context.Context) (PipelineObjectResponse, error) {
	var response PipelineObjectResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	if s.objectName == nil {
		return response, errors.New("hevo: missing required objectName")
	}

	err := s.HTTPService.Do(ctx, http.MethodGet, objectPath(*s.pipelineID, *s.objectName), nil, nil, &response)

	return response, err
}
