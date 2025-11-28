package models

type EnvironmentInputs struct {
	Source  EnvironmentSource
	Version map[string]string  `json:"version"`
	Params  *EnvironmentParams `json:"params,omitempty"`
}

type EnvironmentSource struct {
	Source
	Environment EnvironmentFilter `json:"environment"`
}

type EnvironmentFilter struct {
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Search *string `json:"search,omitempty"`
	States *string `json:"states,omitempty"`
}

type EnvironmentParams struct {
	// Action define what to do on put steps, create, stop or delete.
	Action string `json:"action"`

	// Environments parameters

	// Name of the environment
	Name string `json:"name"`
	// Description of the environment
	Description *string `json:"description"`
	// ExternalURL is the link to this environment
	ExternalURL *string `json:"external_url"`
	// Tier represent the environment tier, must be production, staging, testing, development or other
	Tier *string `json:"tier"`
	// Represent the gitlab ClusterAgentID
	ClusterAgentID *int `json:"cluster_agent_id"`
	// Assign a KubernetesNamespace
	KubernetesNamespace *string `json:"kubernetes_namespace"`
	// Assign a FluxResourcePath
	FluxResourcePath *string `json:"flux_resource_path"`
	// AutoStopSetting must be either always or with_action
	AutoStopSetting *string `json:"auto_stop_setting"`

	// Stop actions parameters

	// For stop action, Force will force stop without executing on_stop action in gitlab
	Force *bool `json:"force,omitempty"`

	// Stop environments that have been modified or deployed to before the specified date.
	// Expected in ISO 8601 format (2019-03-15T08:00:00Z). Valid inputs are between 10 years ago and 1 week ago
	Before *string `json:"before"`
	// For delete action: maximum number of environments to delete. Defaults to 100.
	Limit *string `json:"limit"`
}

// Models related to specific requests
type DeleteStoppedEntries struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	ExternalUrl *string `json:"external_url"`
}
