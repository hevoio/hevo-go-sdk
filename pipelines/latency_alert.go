package pipelines

// LatencyAlert alerts when replication falls behind. Build it with the client's NewLatencyAlert.
type LatencyAlert struct {
	enabled          *bool
	thresholdMinutes *int
}

// LatencyAlertResponse is the Pipeline's latency alert as the API returns it.
type LatencyAlertResponse struct {
	Enabled          bool `json:"enabled"`
	ThresholdMinutes int  `json:"threshold_minutes"`
}

type latencyAlertRequest struct {
	Enabled          *bool `json:"enabled,omitempty"`
	ThresholdMinutes *int  `json:"threshold_minutes,omitempty"`
}

// Enabled turns the alert on or off.
func (a *LatencyAlert) Enabled(value bool) *LatencyAlert {
	a.enabled = &value

	return a
}

// ThresholdMinutes sets the latency that raises the alert.
func (a *LatencyAlert) ThresholdMinutes(value int) *LatencyAlert {
	a.thresholdMinutes = &value

	return a
}

// request returns the alert as the API expects it.
func (a *LatencyAlert) request() *latencyAlertRequest {
	if a == nil {
		return nil
	}

	return &latencyAlertRequest{Enabled: a.enabled, ThresholdMinutes: a.thresholdMinutes}
}
