package mergerequeststatus

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cycloidio/gitlab-resource/models"
)

type Handler struct {
	stderr io.Writer
	stdout io.Writer
	cfg    *models.MergeRequestInputs
}

func NewHandler(stdout, stderr io.Writer, input []byte) (*Handler, error) {
	var config *models.MergeRequestInputs
	err := json.Unmarshal(input, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input config from JSON: %b: %w", input, err)
	}

	return &Handler{
		stdout: stdout,
		stderr: stderr,
		cfg:    config,
	}, nil
}
