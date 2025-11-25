package environments

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

func (h Handler) In(outDir string) error {
	// Get the full user payload to get the user email
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	envID, err := strconv.ParseInt(h.cfg.Version["id"], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to cast current environment id %q: %w", h.cfg.Version["id"], err)
	}

	env, _, err := client.Environments.GetEnvironment(h.cfg.Source.ProjectID, int(envID))
	if err != nil {
		return fmt.Errorf("failed to fetch environment with id %q: %w", h.cfg.Version["id"], err)
	}

	envJSON, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize env to JSON: %w", err)
	}

	metadataPath := path.Join(outDir, "environment.json")
	err = os.WriteFile(metadataPath, envJSON, 0666)
	if err != nil {
		return fmt.Errorf("failed to write version to output dir %q: %w", outDir, err)
	}

	output := models.Output{
		Version:  h.cfg.Version,
		Metadata: EnvironmentToMetadatas(env),
	}

	return internal.OutputJSON(h.stdout, output)
}
