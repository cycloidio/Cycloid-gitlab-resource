package models

type MergeRequestSource struct {
	Source

	// MergeRequestStatus contains the check settings for `merge_request_status` feature
	MergeRequestStatus MergeRequestFilter `json:"merge_request_status"`
}

type MergeRequestState string

const (
	Opened MergeRequestState = "opened"
	Closed MergeRequestState = "closed"
	Locked MergeRequestState = "locked"
	Merged MergeRequestState = "merged"
)

type MergeRequestFilter struct {
	// MergeRequestIID requires the Gitlab merge request ID (not IID)
	MergeRequestIID int `json:"merge_request_iid"`
	// State filter which state will produce a version
	State []MergeRequestState
}

type MergeRequestInputs struct {
	Source  MergeRequestSource
	Version map[string]string   `json:"version"`
	Params  *MergeRequestParams `json:"params,omitempty"`
}

type MergeRequestParams struct {
}
