package models

type Input struct {
	Source Source `json:"source"`
}

type Source struct {
	ProjectID string `json:"project_id"`
	ServerURL string `json:"server_url"`
	Feature   string `json:"feature"`
	Token     string `json:"token"`
}

type Output struct {
	Version  []map[string]string `json:"version"`
	Metadata Metadatas           `json:"metadata"`
}

type Metadatas []Metadata
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
