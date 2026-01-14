package environments

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
)

func (h Handler) In(outDir string) error {
	envIDStr, ok := h.cfg.Version["id"]
	if !ok {
		// In that case, the version is empty, so we output nothing.
		return internal.OutputJSON(h.stdout, models.Output{
			Version:  h.cfg.Version,
			Metadata: models.Metadatas{},
		})
	}

	envID, err := strconv.ParseInt(envIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to cast current environment id %q: %w", h.cfg.Version["id"], err)
	}

	env, _, err := h.glab.Environments.GetEnvironment(h.cfg.Source.ProjectID, int(envID))
	if err != nil {
		return fmt.Errorf("failed to fetch environment with id %q: %w", h.cfg.Version["id"], err)
	}

	envJSON, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize env to JSON: %w", err)
	}

	metadataPath := path.Join(outDir, "metadata.json")
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
