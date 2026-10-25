// Package hevo is the Go SDK for Hevo's API.
//
// Create a Client with New, then build each call from the method named after its operation and
// send it with Do:
//
//	client := hevo.New(accessKey, secretKey)
//	if err := client.Region("us"); err != nil {
//		return err
//	}
//
//	resp, err := client.GetPipeline().PipelineID(42).Do(ctx)
package hevo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/hevoio/hevo-go/destinations"
	httputils "github.com/hevoio/hevo-go/http_utils"
	"github.com/hevoio/hevo-go/jobs"
	"github.com/hevoio/hevo-go/metadata"
	pipelineobjects "github.com/hevoio/hevo-go/pipeline_objects"
	"github.com/hevoio/hevo-go/pipelines"
)

// Version is the SDK's version. Update it on each release.
const Version = "0.1.0"

// DefaultUserAgent identifies the SDK in every request.
const DefaultUserAgent = "hevo-go/" + Version

const (
	defaultMaxRetryAttempts = 5
	defaultMaxConcurrency   = 4
	defaultTimeout          = 60 * time.Second
)

// regionURLs maps each Hevo region to its API host.
var regionURLs = map[string]string{
	"us":   "https://us.hevodata.com",
	"us2":  "https://us2.hevodata.com",
	"eu":   "https://eu.hevodata.com",
	"eu2":  "https://eu2.hevodata.com",
	"in":   "https://in.hevodata.com",
	"asia": "https://asia.hevodata.com",
	"au":   "https://au.hevodata.com",
}

// Client holds the SDK's configuration. Set it up before building calls; it is then safe for
// concurrent use.
type Client struct {
	baseURL          string
	authorization    string
	customUserAgent  string
	httpClient       httputils.HTTPClient
	handleRateLimits bool
	maxRetryAttempts int
	slots            chan struct{}
}

// New returns a Client that authenticates with a team's API access key and secret. It calls the
// us region until Region or BaseURL says otherwise.
func New(accessKey, secretKey string) *Client {
	credentials := base64.StdEncoding.EncodeToString([]byte(accessKey + ":" + secretKey))

	return &Client{
		baseURL:          regionURLs["us"],
		authorization:    "Basic " + credentials,
		httpClient:       &http.Client{Timeout: defaultTimeout},
		handleRateLimits: true,
		maxRetryAttempts: defaultMaxRetryAttempts,
		slots:            make(chan struct{}, defaultMaxConcurrency),
	}
}

// Region sends requests to a Hevo region: one of Regions().
func (c *Client) Region(region string) error {
	url, ok := regionURLs[strings.ToLower(strings.TrimSpace(region))]
	if !ok {
		return fmt.Errorf("hevo: unknown region %q; use one of %s", region, strings.Join(Regions(), ", "))
	}

	c.baseURL = url

	return nil
}

// BaseURL sends requests to baseURL instead of a region's host, for example
// http://localhost:9080 for a local stack.
func (c *Client) BaseURL(baseURL string) {
	c.baseURL = strings.TrimRight(baseURL, "/")
}

// CustomUserAgent is appended to DefaultUserAgent, for example
// "terraform-provider-hevo/1.0.0 terraform/1.9.0".
func (c *Client) CustomUserAgent(customUserAgent string) {
	c.customUserAgent = customUserAgent
}

// SetHTTPClient sends requests through httpClient instead of an *http.Client with a 60 s
// timeout.
func (c *Client) SetHTTPClient(httpClient httputils.HTTPClient) {
	c.httpClient = httpClient
}

// SetHandleRateLimits turns waiting and retrying on a 429 on or off. It is on by default.
func (c *Client) SetHandleRateLimits(handleRateLimits bool) {
	c.handleRateLimits = handleRateLimits
}

// SetMaxRetryAttempts caps the retries of one call. The default is 5.
func (c *Client) SetMaxRetryAttempts(maxRetryAttempts int) {
	c.maxRetryAttempts = maxRetryAttempts
}

// SetMaxConcurrency caps the calls in flight at n; 0 or less removes the cap. The default is 4.
func (c *Client) SetMaxConcurrency(n int) {
	c.slots = nil
	if n > 0 {
		c.slots = make(chan struct{}, n)
	}
}

// Regions lists the regions Region accepts.
func Regions() []string {
	names := make([]string, 0, len(regionURLs))
	for name := range regionURLs {
		names = append(names, name)
	}

	slices.Sort(names)

	return names
}

// HTTPService returns the settings every call is sent with.
func (c *Client) HTTPService() httputils.HTTPService {
	return httputils.HTTPService{
		BaseURL:          c.baseURL,
		CommonHeaders:    c.commonHeaders(),
		Client:           c.httpClient,
		HandleRateLimits: c.handleRateLimits,
		MaxRetryAttempts: c.maxRetryAttempts,
		Slots:            c.slots,
	}
}

func (c *Client) commonHeaders() map[string]string {
	userAgent := DefaultUserAgent
	if c.customUserAgent != "" {
		userAgent += " " + c.customUserAgent
	}

	return map[string]string{
		"Authorization": c.authorization,
		"User-Agent":    userAgent,
		"Accept":        "application/json",
	}
}

// ListPipelines builds a call; see pipelines.ListPipelinesService.
func (c *Client) ListPipelines() *pipelines.ListPipelinesService {
	return &pipelines.ListPipelinesService{HTTPService: c.HTTPService()}
}

// CreatePipeline builds a call; see pipelines.CreatePipelineService.
func (c *Client) CreatePipeline() *pipelines.CreatePipelineService {
	return &pipelines.CreatePipelineService{HTTPService: c.HTTPService()}
}

// GetPipeline builds a call; see pipelines.GetPipelineService.
func (c *Client) GetPipeline() *pipelines.GetPipelineService {
	return &pipelines.GetPipelineService{HTTPService: c.HTTPService()}
}

// UpdatePipeline builds a call; see pipelines.UpdatePipelineService.
func (c *Client) UpdatePipeline() *pipelines.UpdatePipelineService {
	return &pipelines.UpdatePipelineService{HTTPService: c.HTTPService()}
}

// DeletePipeline builds a call; see pipelines.DeletePipelineService.
func (c *Client) DeletePipeline() *pipelines.DeletePipelineService {
	return &pipelines.DeletePipelineService{HTTPService: c.HTTPService()}
}

// EnablePipeline builds a call; see pipelines.EnablePipelineService.
func (c *Client) EnablePipeline() *pipelines.EnablePipelineService {
	return &pipelines.EnablePipelineService{HTTPService: c.HTTPService()}
}

// DisablePipeline builds a call; see pipelines.DisablePipelineService.
func (c *Client) DisablePipeline() *pipelines.DisablePipelineService {
	return &pipelines.DisablePipelineService{HTTPService: c.HTTPService()}
}

// SyncPipelineNow builds a call; see pipelines.SyncPipelineNowService.
func (c *Client) SyncPipelineNow() *pipelines.SyncPipelineNowService {
	return &pipelines.SyncPipelineNowService{HTTPService: c.HTTPService()}
}

// ResyncPipeline builds a call; see pipelines.ResyncPipelineService.
func (c *Client) ResyncPipeline() *pipelines.ResyncPipelineService {
	return &pipelines.ResyncPipelineService{HTTPService: c.HTTPService()}
}

// UpdatePipelineSource builds a call; see pipelines.UpdatePipelineSourceService.
func (c *Client) UpdatePipelineSource() *pipelines.UpdatePipelineSourceService {
	return &pipelines.UpdatePipelineSourceService{HTTPService: c.HTTPService()}
}

// GetPipelinePosition builds a call; see pipelines.GetPipelinePositionService.
func (c *Client) GetPipelinePosition() *pipelines.GetPipelinePositionService {
	return &pipelines.GetPipelinePositionService{HTTPService: c.HTTPService()}
}

// ListPipelineObjects builds a call; see pipelineobjects.ListPipelineObjectsService.
func (c *Client) ListPipelineObjects() *pipelineobjects.ListPipelineObjectsService {
	return &pipelineobjects.ListPipelineObjectsService{HTTPService: c.HTTPService()}
}

// GetPipelineObject builds a call; see pipelineobjects.GetPipelineObjectService.
func (c *Client) GetPipelineObject() *pipelineobjects.GetPipelineObjectService {
	return &pipelineobjects.GetPipelineObjectService{HTTPService: c.HTTPService()}
}

// UpdatePipelineObjects builds a call; see pipelineobjects.UpdatePipelineObjectsService.
func (c *Client) UpdatePipelineObjects() *pipelineobjects.UpdatePipelineObjectsService {
	return &pipelineobjects.UpdatePipelineObjectsService{HTTPService: c.HTTPService()}
}

// ResyncPipelineObjects builds a call; see pipelineobjects.ResyncPipelineObjectsService.
func (c *Client) ResyncPipelineObjects() *pipelineobjects.ResyncPipelineObjectsService {
	return &pipelineobjects.ResyncPipelineObjectsService{HTTPService: c.HTTPService()}
}

// ListPipelineJobs builds a call; see jobs.ListPipelineJobsService.
func (c *Client) ListPipelineJobs() *jobs.ListPipelineJobsService {
	return &jobs.ListPipelineJobsService{HTTPService: c.HTTPService()}
}

// GetPipelineJob builds a call; see jobs.GetPipelineJobService.
func (c *Client) GetPipelineJob() *jobs.GetPipelineJobService {
	return &jobs.GetPipelineJobService{HTTPService: c.HTTPService()}
}

// ListPipelineJobObjects builds a call; see jobs.ListPipelineJobObjectsService.
func (c *Client) ListPipelineJobObjects() *jobs.ListPipelineJobObjectsService {
	return &jobs.ListPipelineJobObjectsService{HTTPService: c.HTTPService()}
}

// CancelPipelineJob builds a call; see jobs.CancelPipelineJobService.
func (c *Client) CancelPipelineJob() *jobs.CancelPipelineJobService {
	return &jobs.CancelPipelineJobService{HTTPService: c.HTTPService()}
}

// ListDestinations builds a call; see destinations.ListDestinationsService.
func (c *Client) ListDestinations() *destinations.ListDestinationsService {
	return &destinations.ListDestinationsService{HTTPService: c.HTTPService()}
}

// CreateDestination builds a call; see destinations.CreateDestinationService.
func (c *Client) CreateDestination() *destinations.CreateDestinationService {
	return &destinations.CreateDestinationService{HTTPService: c.HTTPService()}
}

// GetDestination builds a call; see destinations.GetDestinationService.
func (c *Client) GetDestination() *destinations.GetDestinationService {
	return &destinations.GetDestinationService{HTTPService: c.HTTPService()}
}

// UpdateDestination builds a call; see destinations.UpdateDestinationService.
func (c *Client) UpdateDestination() *destinations.UpdateDestinationService {
	return &destinations.UpdateDestinationService{HTTPService: c.HTTPService()}
}

// DeleteDestination builds a call; see destinations.DeleteDestinationService.
func (c *Client) DeleteDestination() *destinations.DeleteDestinationService {
	return &destinations.DeleteDestinationService{HTTPService: c.HTTPService()}
}

// ListSourceConnectors builds a call; see metadata.ListSourceConnectorsService.
func (c *Client) ListSourceConnectors() *metadata.ListSourceConnectorsService {
	return &metadata.ListSourceConnectorsService{HTTPService: c.HTTPService()}
}

// GetSourceConnector builds a call; see metadata.GetSourceConnectorService.
func (c *Client) GetSourceConnector() *metadata.GetSourceConnectorService {
	return &metadata.GetSourceConnectorService{HTTPService: c.HTTPService()}
}

// ListDestinationConnectors builds a call; see metadata.ListDestinationConnectorsService.
func (c *Client) ListDestinationConnectors() *metadata.ListDestinationConnectorsService {
	return &metadata.ListDestinationConnectorsService{HTTPService: c.HTTPService()}
}

// GetDestinationConnector builds a call; see metadata.GetDestinationConnectorService.
func (c *Client) GetDestinationConnector() *metadata.GetDestinationConnectorService {
	return &metadata.GetDestinationConnectorService{HTTPService: c.HTTPService()}
}
