package deployments

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h Handler) Out(outDir string) error {
	// Out script has the Version in the metadata.json
	entries, err := os.ReadDir(outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
	}

	for _, entry := range entries {
		fmt.Fprintf(os.Stderr, "%s\n", entry.Name())
	}

	metadataPath := outDir + "/metadata.json"
	versionBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read current metadata at %q: %w", metadataPath, err)
	}

	err = json.Unmarshal(versionBytes, &h.cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to read the version from %q: %w", metadataPath, err)
	}

	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	switch h.cfg.Params.Action {
	case "create":
		deploy, _, err := client.Deployments.CreateProjectDeployment(
			h.cfg.Source.ProjectID, &gitlab.CreateProjectDeploymentOptions{
				Environment: h.cfg.Source.Environment,
				SHA:         h.cfg.Params.SHA,
				Ref:         h.cfg.Params.Ref,
				Tag:         h.cfg.Params.Tag,
				Status:      (*gitlab.DeploymentStatusValue)(&h.cfg.Params.Status),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create deployment: %w", err)
		}

		deployJSON, err := json.MarshalIndent(deploy, "", "  ")
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			_, _ = fmt.Fprintf(h.stderr, "failed to serialize deployment payload to JSON: %s\n", err.Error())
		}

		err = os.WriteFile(outDir+"/metadata.json", deployJSON, 0666)
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			_, _ = fmt.Fprintf(h.stderr, "failed to write metadata to %q: %s\n", outDir+"/metadata.json", err.Error())
		}

		metadata := models.Metadatas{
			{Name: "id", Value: strconv.FormatInt(int64(deploy.ID), 10)},
			{Name: "status", Value: deploy.Status},
			{Name: "ref", Value: deploy.Ref},
			{Name: "sha", Value: deploy.SHA},
			{Name: "deployable_name", Value: deploy.Deployable.Name},
			{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(deploy.Deployable.Pipeline.ID), 10)},
			{Name: "environment_id", Value: strconv.FormatInt(int64(deploy.Environment.ID), 10)},
			{Name: "environment_name", Value: deploy.Environment.Name},
			{Name: "environment_external_url", Value: deploy.Environment.ExternalURL},
		}

		if deploy.Deployable.Commit != nil {
			metadata = append(metadata, models.Metadatas{
				{Name: "deployable_commit_author_name", Value: deploy.Deployable.Commit.AuthorName},
				{Name: "deployable_commit_author_email", Value: deploy.Deployable.Commit.AuthorEmail},
			}...)
		}

		output := &models.Output{
			Version:  DeploymentToVersion(deploy),
			Metadata: metadata,
		}

		return OutputJSON(h.stdout, output)

	case "update":
		deployIDStr, ok := h.cfg.Version["id"]
		if !ok {
			return fmt.Errorf("failed to update deployment, missing deployment ID in version %v", h.cfg.Version)
		}

		deployID, err := strconv.ParseInt(deployIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse deployment id %q: %w", deployIDStr, err)
		}

		updatedDeploy, _, err := client.Deployments.UpdateProjectDeployment(
			h.cfg.Source.ProjectID, int(deployID), &gitlab.UpdateProjectDeploymentOptions{
				Status: (*gitlab.DeploymentStatusValue)(&h.cfg.Params.Status),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update deployment with id %d: %w", deployID, err)
		}
		userIDStr, ok := h.cfg.Version["user_id"]
		if !ok {
			return fmt.Errorf("cannot get user id from version %v", h.cfg.Version)
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse user id %q: %w", userIDStr, err)
		}

		user, err := internal.GetUser(int(userID), client)
		if err != nil {
			return err
		}

		deployJSON, err := json.MarshalIndent(updatedDeploy, "", "  ")
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			fmt.Fprintf(os.Stderr, "failed to serialize deployment payload to JSON: %s\n", err.Error())
		}

		err = os.WriteFile(outDir+"/metadata.json", deployJSON, 0666)
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			fmt.Fprintf(os.Stderr, "failed to write metadata to %q: %s\n", outDir+"/metadata.json", err.Error())
		}

		metadata := models.Metadatas{
			{Name: "id", Value: strconv.FormatInt(int64(updatedDeploy.ID), 10)},
			{Name: "status", Value: updatedDeploy.Status},
			{Name: "ref", Value: updatedDeploy.Ref},
			{Name: "sha", Value: updatedDeploy.SHA},
			{Name: "deployable_name", Value: updatedDeploy.Deployable.Name},
			{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(updatedDeploy.Deployable.Pipeline.ID), 10)},
			{Name: "environment_id", Value: strconv.FormatInt(int64(updatedDeploy.Environment.ID), 10)},
			{Name: "environment_name", Value: updatedDeploy.Environment.Name},
			{Name: "environment_external_url", Value: updatedDeploy.Environment.ExternalURL},
			{Name: "user_username", Value: user.Username},
			{Name: "user_email", Value: user.Email},
		}

		if updatedDeploy.Deployable.Commit != nil {
			metadata = append(metadata, models.Metadatas{
				{Name: "deployable_commit_author_name", Value: updatedDeploy.Deployable.Commit.AuthorName},
				{Name: "deployable_commit_author_email", Value: updatedDeploy.Deployable.Commit.AuthorEmail},
			}...)
		}

		output := &models.Output{
			Version:  DeploymentToVersion(updatedDeploy),
			Metadata: metadata,
		}
		return OutputJSON(h.stdout, output)

	case "delete":
		deployIDStr, ok := h.cfg.Version["id"]
		if !ok {
			return fmt.Errorf("failed to update deployment, missing deployment ID in version %v", h.cfg.Version)
		}

		deployID, err := strconv.ParseInt(deployIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse deployment id %q: %w", deployIDStr, err)
		}

		_, err = client.Deployments.DeleteProjectDeployment(h.cfg.Source.ProjectID, int(deployID), nil)
		if err != nil {
			return fmt.Errorf("failed to update deployment with id %q: %w", deployIDStr, err)
		}

		output := &models.Output{
			Version:  nil,
			Metadata: models.Metadatas{},
		}
		return OutputJSON(h.stdout, output)

	default:
		return fmt.Errorf("invalid params.action parameter, accepted values are: create, update, delete")
	}
}
