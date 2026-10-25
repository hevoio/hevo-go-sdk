package pipelineobjects

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
	"github.com/hevoio/hevo-go/pipelines"
)

// ResyncPipelineObjectsService reloads some objects from the start. It runs asynchronously.
type ResyncPipelineObjectsService struct {
	httputils.HTTPService

	pipelineID  *int64
	objectNames []string
	dropAndLoad *bool
}

// PipelineID sets the Pipeline.
func (s *ResyncPipelineObjectsService) PipelineID(value int64) *ResyncPipelineObjectsService {
	s.pipelineID = &value

	return s
}

// ObjectNames sets the objects to reload.
func (s *ResyncPipelineObjectsService) ObjectNames(values ...string) *ResyncPipelineObjectsService {
	s.objectNames = values

	return s
}

// DropAndLoad drops the Destination tables before reloading them.
func (s *ResyncPipelineObjectsService) DropAndLoad(value bool) *ResyncPipelineObjectsService {
	s.dropAndLoad = &value

	return s
}

// Do starts the resync.
func (s *ResyncPipelineObjectsService) Do(ctx context.Context) (pipelines.PipelineDetailsResponse, error) {
	var response pipelines.PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	if len(s.objectNames) == 0 {
		return response, errors.New("hevo: missing required objectNames")
	}

	req := &pipelineObjectsResyncRequest{ObjectNames: s.objectNames, DropAndLoad: s.dropAndLoad}
	err := s.HTTPService.Do(ctx, http.MethodPost, objectsPath(*s.pipelineID, "/actions/resync"), req, nil, &response)

	return response, err
}
