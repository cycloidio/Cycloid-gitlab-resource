package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
)

func (h *Handler) In(outDir string) error {
	metadataPath := path.Join(outDir, "metadata.json")

	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	pipelineID, err := strconv.ParseInt(h.cfg.Version["id"], 10, 32)
	if err != nil {
		return fmt.Errorf("missing id in version or bad id from %v: %w", h.cfg.Version, err)
	}

	pipeline, _, err := glab.Pipelines.GetPipeline(h.cfg.Source.ProjectID, int(pipelineID))
	if err != nil {
		return fmt.Errorf(
			"failed to get pipeline with id %q from project id %q: %w",
			h.cfg.Version["id"], h.cfg.Source.ProjectID, err,
		)
	}

	output := &models.Output{
		Version:  PipelinetoVersion(pipeline),
		Metadata: PipelinetoMetadatas(pipeline),
	}

	pipelineJSON, err := json.Marshal(h.cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to serialize pipeline response to JSON: %w", err)
	}

	err = os.WriteFile(metadataPath, pipelineJSON, 0666)
	if err != nil {
		return fmt.Errorf("failed to write pipeline to %q: %w", metadataPath, err)
	}

	return internal.OutputJSON(h.stdout, output)
}
