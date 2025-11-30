package environments

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func EnvironmentToVersion(env *gitlab.Environment) map[string]string {
	var version = make(map[string]string)
	version["id"] = strconv.FormatInt(int64(env.ID), 10)
	version["name"] = env.Name
	version["slug"] = env.Slug
	version["description"] = env.Description
	version["external_url"] = env.ExternalURL
	version["state"] = env.State
	version["tier"] = env.Tier
	version["kubernetes_namespace"] = env.KubernetesNamespace
	version["flux_resource_path"] = env.FluxResourcePath
	version["auto_stop_settin"] = env.AutoStopSetting
	if env.AutoStopAt != nil {
		version["auto_stop_at"] = env.AutoStopAt.String()
	}

	if env.CreatedAt != nil {
		version["created_at"] = env.CreatedAt.String()
	}

	if env.UpdatedAt != nil {
		version["updated_at"] = env.UpdatedAt.String()
	}

	if env.LastDeployment != nil {
		version["last_deployment_id"] = strconv.FormatInt(int64(env.LastDeployment.ID), 10)
		version["last_deployment_status"] = env.LastDeployment.Status
		version["last_deployment_sha"] = env.LastDeployment.SHA
	}
	return version
}

func EnvironmentsToVersion(environments []*gitlab.Environment) []map[string]string {
	var version = make([]map[string]string, len(environments))
	for i, env := range environments {
		version[i] = EnvironmentToVersion(env)
	}

	return version
}

func EnvironmentToMetadatas(env *gitlab.Environment) models.Metadatas {
	return models.Metadatas{
		{Name: "id", Value: strconv.Itoa(env.ID)},
		{Name: "name", Value: env.Name},
		{Name: "slug", Value: env.Slug},
		{Name: "description", Value: env.Description},
		{Name: "state", Value: env.State},
		{Name: "external_url", Value: env.ExternalURL},
	}
}

// ReadDataFromFile will read the current metadata JSON contained in metadataDir
func ReadDataFromFile(metadataDir string) (*gitlab.Environment, error) {
	metaFile := path.Join(metadataDir, "metadata.json")
	var currentMetadata *gitlab.Environment
	if _, err := os.Stat(metaFile); err == nil {
		content, err := os.ReadFile(metaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read current metadata file %q: %w", metaFile, err)
		}

		err = json.Unmarshal(content, &currentMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decode current data as JSON from %q: %w", metaFile, err)
		}
	} else {
		return nil, fmt.Errorf("failed to read current metadata, file not found at %q", metaFile)
	}

	return currentMetadata, nil
}
