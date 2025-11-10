package deployments

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cycloidio/gitlab-resource/models"
)

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

	return &Handler{
		stdout: stdout,
		stderr: stderr,
		cfg:    config,
	}, nil
}

type Handler struct {
	stdout io.Writer
	stderr io.Writer
	cfg    *models.DeploymentInputs
}
