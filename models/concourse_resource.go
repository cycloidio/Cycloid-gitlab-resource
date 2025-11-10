package models

type Inputs struct {
	Source  Source  `json:"source"`
	Version any     `json:"version"`
	Params  *Params `json:"params,omitempty"`
}

type Params struct {
	Deployments ParamDeployments `json:"deployments"`
}

type ParamDeployments struct {
	Action string       `json:"action"`
	Create *ParamCreate `json:"create,omitempty"`
	Update *ParamUpdate `json:"update,omitempty"`
}

type ParamCreate struct {
	SHA    string `json:"sha"`
	Ref    string `json:"ref"`
	Tag    bool   `json:"tag"`
	Status string `json:"status"`
}

type ParamUpdate struct {
	Status string `json:"status"`
}

type Source struct {
	ProjectID   string               `json:"project_id"`
	ServerURL   string               `json:"server_url"`
	Feature     string               `json:"feature"`
	Deployments *SourceDeployment    `json:"deployments"`
	Auth        SourceAuthentication `json:"auth"`
}

type SourceDeployment struct {
	Status      *string `json:"status,omitempty"`
	Environment *string `json:"environment,omitempty"`
}

type SourceAuthentication struct {
	Token string `json:"token"`
}

type Output struct {
	Version  any       `json:"version"`
	Metadata Metadatas `json:"metadata"`
}

type Metadatas []Metadata
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BuildMetadata struct{}
