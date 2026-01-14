package deployments

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/models"
)

var (
	AvailableModes = []string{"classic", "yield"}
)

type Handler struct {
	stdout io.Writer
	stderr io.Writer
	cfg    *models.DeploymentInputs
	glab   *gitlab.Client
}

func NewHandler(stdout, stderr io.Writer, input []byte) (*Handler, error) {
	var config *models.DeploymentInputs
	err := json.Unmarshal(input, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode source from stdin: %s: %w", input, err)
	}

	if config.Source.Token == "" {
		return nil, fmt.Errorf("auth token empty, please provide a gitlab token")
	}

	if config.Source.ServerURL == "" {
		return nil, fmt.Errorf("missing server_url parameter")
	}

	if !slices.Contains(AvailableModes, config.Source.Mode) {
		return nil, fmt.Errorf("source.mode must one of those values: %s", strings.Join(AvailableModes, ", "))
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
