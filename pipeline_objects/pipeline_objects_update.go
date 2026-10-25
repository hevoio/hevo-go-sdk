package pipelineobjects

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
	"github.com/hevoio/hevo-go/pipelines"
)

// UpdatePipelineObjectsService changes which objects a Pipeline replicates. SelectionMode says
// what happens to the objects not named: RETAIN_EXISTING (the default) leaves them,
// INCLUDE_ALL selects them and EXCLUDE_ALL deselects them. So EXCLUDE_ALL with the wanted
// objects selected sets exactly that set. It runs asynchronously, like a Pipeline update.
type UpdatePipelineObjectsService struct {
	httputils.HTTPService

	pipelineID    *int64
	selectionMode *string
	objects       []*ObjectSelection
}

// PipelineID sets the Pipeline.
func (s *UpdatePipelineObjectsService) PipelineID(value int64) *UpdatePipelineObjectsService {
	s.pipelineID = &value

	return s
}

// SelectionMode sets INCLUDE_ALL, EXCLUDE_ALL or RETAIN_EXISTING.
func (s *UpdatePipelineObjectsService) SelectionMode(value string) *UpdatePipelineObjectsService {
	s.selectionMode = &value

	return s
}

// Object adds one object to change.
func (s *UpdatePipelineObjectsService) Object(value *ObjectSelection) *UpdatePipelineObjectsService {
	s.objects = append(s.objects, value)

	return s
}

// Do updates the selection.
func (s *UpdatePipelineObjectsService) Do(ctx context.Context) (pipelines.PipelineDetailsResponse, error) {
	var response pipelines.PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	req := &pipelineObjectsUpdateRequest{SelectionMode: s.selectionMode}
	for _, object := range s.objects {
		req.Objects = append(req.Objects, object.request())
	}

	err := s.HTTPService.Do(ctx, http.MethodPatch, objectsPath(*s.pipelineID, ""), req, nil, &response)

	return response, err
}
