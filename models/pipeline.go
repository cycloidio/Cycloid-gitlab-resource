package models

import gitlab "gitlab.com/gitlab-org/api/client-go"

type PipelineSource struct {
	Source

	Pipeline PipelineFilter `json:"pipeline"`
}

type CheckMode string

const (
	// List      CheckMode = "list"
	Get CheckMode = "get"
)

type PipelineFilter struct {
	CheckMode *CheckMode `json:"check_mode"`

	// Ref to get when using get mode
	Ref *string `json:"ref"`
}

type PipelineInputs struct {
	Source  PipelineSource
	Version map[string]string `json:"version"`
	Params  *PipelineParams   `json:"params"`
}

type PipelineAction string

const (
	ActionCreate PipelineAction = "create"
)

type PipelineParams struct {
	Action PipelineAction

	// For create action

	// Either Ref or MergeRequest is required
	// Ref from git
	Ref *string `json:"ref,omitempty"`
	// Merge request IID
	MergeRequestIID *int `json:"merge_request_iid,omitempty"`

	Variables *[]*gitlab.PipelineVariableOptions
	Inputs    map[string]any
}
