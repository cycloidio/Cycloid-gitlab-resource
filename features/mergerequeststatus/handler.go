package mergerequeststatus

import (
	"encoding/json"
	"fmt"
	"io"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Handler struct {
	stderr io.Writer
	stdout io.Writer
	cfg    *models.MergeRequestInputs
	glab   *gitlab.Client
}

func NewHandler(stdout, stderr io.Writer, input []byte) (*Handler, error) {
	var config *models.MergeRequestInputs
	err := json.Unmarshal(input, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input config from JSON: %b: %w", input, err)
	}

	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: config.Source.Token,
		Url:   config.Source.ServerURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gitlab client: %w", err)
	}

	return &Handler{
		stdout: stdout,
		stderr: stderr,
		cfg:    config,
		glab:   glab,
	}, nil
}
