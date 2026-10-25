package pipelines

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/hevoio/hevo-go/http_utils"
)

// EnablePipelineService starts a Pipeline replicating. It needs at least one selected object. It runs asynchronously: the Pipeline reports Enabling, then Initialized.
type EnablePipelineService struct {
	httputils.HTTPService

	pipelineID *int64
}

// PipelineID sets the Pipeline.
func (s *EnablePipelineService) PipelineID(value int64) *EnablePipelineService {
	s.pipelineID = &value

	return s
}

// Do sends the request.
func (s *EnablePipelineService) Do(ctx context.Context) (PipelineDetailsResponse, error) {
	var response PipelineDetailsResponse

	if s.pipelineID == nil {
		return response, errors.New("hevo: missing required pipelineID")
	}

	err := s.HTTPService.Do(ctx, http.MethodPost, pipelinePath(*s.pipelineID, "/actions/enable"), nil, nil, &response)

	return response, err
}
