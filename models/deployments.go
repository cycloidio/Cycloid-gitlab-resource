package models

import gitlab "gitlab.com/gitlab-org/api/client-go"

type DeploymentInputs struct {
	Source  DeploymentSource
	Version *gitlab.Deployment
	Params  *ParamDeployments `json:"params,omitempty"`
}

type DeploymentSource struct {
	Source
	Status      *string `json:"status,omitempty"`
	Environment *string `json:"environment,omitempty"`
}

type ParamDeployments struct {
	Action string  `json:"action"`
	SHA    *string `json:"sha,omitempty"`
	Ref    *string `json:"ref,omitempty"`
	Tag    *bool   `json:"tag,omitempty"`
	Status string  `json:"status"`
}
