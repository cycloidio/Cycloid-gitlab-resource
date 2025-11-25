package environments

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h Handler) Out(outDir string) error {
	if h.cfg.Params.MetadataDir == nil {
		return fmt.Errorf("missing metadata_dir parameter for PUT")
	}
	metadataDir := path.Join(outDir, *h.cfg.Params.MetadataDir)

	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	var output *models.Output
	switch h.cfg.Params.Action {
	case "create":
		env, _, err := client.Environments.CreateEnvironment(h.cfg.Source.ProjectID, &gitlab.CreateEnvironmentOptions{
			Name:                &h.cfg.Params.Name,
			Description:         h.cfg.Params.Description,
			ExternalURL:         h.cfg.Params.ExternalURL,
			Tier:                h.cfg.Params.Tier,
			ClusterAgentID:      h.cfg.Params.ClusterAgentID,
			KubernetesNamespace: h.cfg.Params.KubernetesNamespace,
			FluxResourcePath:    h.cfg.Params.FluxResourcePath,
			AutoStopSetting:     h.cfg.Params.AutoStopSetting,
		})
		if err != nil {
			return fmt.Errorf("failed to create environment %q: %w", h.cfg.Params.Name, err)
		}

		output = &models.Output{
			Version:  EnvironmentToVersion(env),
			Metadata: EnvironmentToMetadatas(env),
		}

	case "update":
		env, _, err := client.Environments.EditEnvironment(
			h.cfg.Source.ProjectID, *h.cfg.Params.ID,
			&gitlab.EditEnvironmentOptions{
				Name:                &h.cfg.Params.Name,
				Description:         h.cfg.Params.Description,
				Tier:                h.cfg.Params.Tier,
				AutoStopSetting:     h.cfg.Params.AutoStopSetting,
				ClusterAgentID:      h.cfg.Params.ClusterAgentID,
				KubernetesNamespace: h.cfg.Params.KubernetesNamespace,
				ExternalURL:         h.cfg.Params.ExternalURL,
				FluxResourcePath:    h.cfg.Params.FluxResourcePath,
			})
		if err != nil {
			return fmt.Errorf("failed to update env %q: %w", h.cfg.Params.Name, err)
		}

		output = &models.Output{
			Version:  EnvironmentToVersion(env),
			Metadata: EnvironmentToMetadatas(env),
		}
	case "delete":
		_, err := client.Environments.DeleteEnvironment(h.cfg.Source.ProjectID, *h.cfg.Params.ID)
		if err != nil {
			return fmt.Errorf("failed to delete environment %q: %w", h.cfg.Params.Name, err)
		}

		output = &models.Output{
			Version: map[string]string{},
			Metadata: models.Metadatas{
				{Name: "id", Value: strconv.Itoa(int(*h.cfg.Params.ID))},
				{Name: "state", Value: "deleted"},
			},
		}
	case "delete_stopped":
		url := fmt.Sprintf("projects/%s/environments/review_apps", h.cfg.Source.ProjectID)
		payload := struct {
			Before  *string `json:"before,omitempty"`
			Limit   *string `json:"limit,omitempty"`
			Dry_run *bool   `json:"dry_run,omitempty"`
		}{
			Before:  h.cfg.Params.Before,
			Limit:   h.cfg.Params.Limit,
			Dry_run: gitlab.Ptr(false),
		}

		req, err := client.NewRequest(http.MethodDelete, url, payload, nil)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		response := struct {
			ScheduledEntries     []models.DeleteStoppedEntries `json:"scheduled_entries"`
			UnprocessableEntries []models.DeleteStoppedEntries `json:"unprocessable_entries"`
		}{}
		_, err = client.Do(req, &response)
		if err != nil {
			return fmt.Errorf("failed to delete stopped environments: %w", err)
		}

		var metadatas = make(models.Metadatas, len(response.ScheduledEntries))
		for i, v := range response.ScheduledEntries {
			metadatas[i] = models.Metadata{Name: strconv.Itoa(int(v.Id)), Value: fmt.Sprintf("deleted %q", v.Name)}
		}

		output = &models.Output{
			Version:  map[string]string{},
			Metadata: metadatas,
		}
	case "stop":
		env, _, err := client.Environments.StopEnvironment(
			h.cfg.Source.ProjectID, *h.cfg.Params.ID,
			&gitlab.StopEnvironmentOptions{
				Force: h.cfg.Params.Force,
			})
		if err != nil {
			return fmt.Errorf("failed to stop env %q: %w", h.cfg.Params.Name, err)
		}

		output = &models.Output{
			Version:  EnvironmentToVersion(env),
			Metadata: EnvironmentToMetadatas(env),
		}
	case "stop_stale":
		url := fmt.Sprintf("projects/%s/environments/stop_stale", h.cfg.Source.ProjectID)
		payload := struct {
			Before *string `json:"before,omitempty"`
		}{
			Before: h.cfg.Params.Before,
		}

		req, err := client.NewRequest(http.MethodDelete, url, payload, nil)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		response := struct {
			Message string `json:"message"`
		}{}
		_, err = client.Do(req, &response)
		if err != nil {
			return fmt.Errorf("failed to delete stopped environments: %w", err)
		}

		output = &models.Output{
			Version: map[string]string{},
			Metadata: models.Metadatas{
				{Name: "message", Value: response.Message},
			},
		}
	default:
		return fmt.Errorf("invalid params.action parameter, accepted values are: create, update, delete, delete_stopped, stop, stop_stale")
	}

	err = internal.WriteMetadata(metadataDir, output.Version)
	if err != nil {
		return err
	}

	return internal.OutputJSON(h.stdout, output)
}
