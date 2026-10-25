package pipelines

// Schedule says when the Pipeline ingests. Build it with the client's NewSchedule.
type Schedule struct {
	syncType         *string
	frequencyMinutes *int64
}

// ScheduleResponse is the Pipeline's schedule as the API returns it.
type ScheduleResponse struct {
	SyncType         string `json:"sync_type"`
	FrequencyMinutes int64  `json:"frequency_minutes"`
}

type scheduleRequest struct {
	SyncType         *string `json:"sync_type,omitempty"`
	FrequencyMinutes *int64  `json:"frequency_minutes,omitempty"`
}

// Schedule sync types.
const (
	SyncTypeScheduled = "SCHEDULED"
	SyncTypeOnDemand  = "ON_DEMAND"
)

// SyncType sets SCHEDULED or ON_DEMAND.
func (s *Schedule) SyncType(value string) *Schedule {
	s.syncType = &value

	return s
}

// FrequencyMinutes sets how often a SCHEDULED Pipeline ingests.
func (s *Schedule) FrequencyMinutes(value int64) *Schedule {
	s.frequencyMinutes = &value

	return s
}

// request returns the schedule as the API expects it.
func (s *Schedule) request() *scheduleRequest {
	if s == nil {
		return nil
	}

	return &scheduleRequest{SyncType: s.syncType, FrequencyMinutes: s.frequencyMinutes}
}
