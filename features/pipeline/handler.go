package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Handler struct {
	stderr io.Writer
	stdout io.Writer
	cfg    *models.PipelineInputs
	glab   *gitlab.Client
	logger *slog.Logger
}

func NewHandler(stdout, stderr io.Writer, input []byte) (*Handler, error) {
	var config *models.PipelineInputs
	err := json.Unmarshal(input, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input config from JSON:\n%s %w", string(input), err)
	}

	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: config.Source.Token,
		Url:   config.Source.ServerURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gitlab client: %w", err)
	}

	logOpt := &slog.HandlerOptions{}
	if config.Source.LogLevel == nil {
		logOpt.Level = slog.LevelInfo
	} else {
		switch *config.Source.LogLevel {
		case "debug":
			logOpt.Level = slog.LevelDebug
		case "info":
			logOpt.Level = slog.LevelInfo
		case "warning":
			logOpt.Level = slog.LevelWarn
		case "error":
			logOpt.Level = slog.LevelError
		}
	}

	// we log on stdou since cc don't display stderr traces
	logger := slog.New(slog.NewTextHandler(stdout, logOpt))

	return &Handler{
		stdout: stdout,
		stderr: stderr,
		cfg:    config,
		glab:   glab,
		logger: logger,
	}, nil
}
