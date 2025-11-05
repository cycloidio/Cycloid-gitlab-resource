package models

type CheckInput struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Version []map[string]string

type InInputs struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

type InParams struct {
}

type OutInputs struct {
	Source  Source    `json:"source"`
	Version Version   `json:"version"`
	Params  OutParams `json:"params"`
}

type OutParams struct{}

type Source struct {
	Project     string               `json:"project"`
	ProjectID   string               `json:"project_id"`
	Environment string               `json:"environment"`
	ServerURL   string               `json:"server_url"`
	Feature     string               `json:"feature"`
	Auth        SourceAuthentication `json:"auth"`
}

type SourceAuthentication struct {
	Token string `json:"token"`
}

type Output struct {
	Version  Version   `json:"version"`
	Metadata Metadatas `json:"metadata"`
}

type Metadatas []Metadata
type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BuildMetadata struct{}
