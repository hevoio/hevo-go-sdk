// Package common holds the response parts every Hevo API resource shares.
package common

// CommonResponse is the envelope of every response: Code is "Success", or "Accepted" when the
// work runs asynchronously.
type CommonResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ListResponse carries the cursor of a paginated list. Send NextCursor back as the next request's
// cursor while HasMore is true.
type ListResponse struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// Response codes the API returns in CommonResponse.Code.
const (
	CodeSuccess  = "Success"
	CodeAccepted = "Accepted"
)
