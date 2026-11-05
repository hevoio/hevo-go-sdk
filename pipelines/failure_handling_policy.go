package pipelines

// FailureHandlingPolicy says how many objects may fail before a load fails. Build it with the
// client's NewFailureHandlingPolicy.
type FailureHandlingPolicy struct {
	objectFailureLevel     *string
	objectFailureThreshold *int64
}

// FailureHandlingPolicyResponse is the Pipeline's policy as the API returns it.
type FailureHandlingPolicyResponse struct {
	ObjectFailureLevel     string `json:"object_failure_level"`
	ObjectFailureThreshold int64  `json:"object_failure_threshold"`
}

type failureHandlingPolicyRequest struct {
	ObjectFailureLevel     *string `json:"object_failure_level,omitempty"`
	ObjectFailureThreshold *int64  `json:"object_failure_threshold,omitempty"`
}

// Object failure levels.
const (
	ObjectFailureLevelPartial   = "PARTIAL"
	ObjectFailureLevelAllOrNone = "ALL_OR_NONE"
)

// ObjectFailureLevel sets PARTIAL or ALL_OR_NONE.
func (p *FailureHandlingPolicy) ObjectFailureLevel(value string) *FailureHandlingPolicy {
	p.objectFailureLevel = &value

	return p
}

// ObjectFailureThreshold sets how many objects may fail under PARTIAL.
func (p *FailureHandlingPolicy) ObjectFailureThreshold(value int64) *FailureHandlingPolicy {
	p.objectFailureThreshold = &value

	return p
}

// request returns the policy as the API expects it.
func (p *FailureHandlingPolicy) request() *failureHandlingPolicyRequest {
	if p == nil {
		return nil
	}

	return &failureHandlingPolicyRequest{
		ObjectFailureLevel:     p.objectFailureLevel,
		ObjectFailureThreshold: p.objectFailureThreshold,
	}
}
