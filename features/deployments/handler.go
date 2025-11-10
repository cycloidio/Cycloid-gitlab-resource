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

func NewHandler() Handler {
	return Handler{}
}

type Handler struct{}

func (h Handler) Validate(input *models.Inputs) error {
	if input.Source.Auth.Token == "" {
		return fmt.Errorf("auth token empty, please provide a gitlab token")
	}

	if input.Source.ServerURL == "" {
		return fmt.Errorf("missing server_url parameter")
	}

	if input.Source.Deployments == nil {
		return fmt.Errorf("when using feature 'deployments': 'source.deployments' source parameter cannot be null, please check documentation, current config: %v", internal.MustJSON(input))
	}

	return nil
}

func (h Handler) Check(input *models.Inputs) (any, error) {
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: input.Source.Auth.Token,
		Url:   input.Source.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	var options = &gitlab.ListProjectDeploymentsOptions{
		Sort:        gitlab.Ptr("desc"), // version should be oldest first
		Environment: input.Source.Deployments.Environment,
		Status:      input.Source.Deployments.Status,
	}

	deployments, _, err := client.Deployments.ListProjectDeployments(input.Source.ProjectID, options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments from gitlab API: %w", err)
	}

	return GetVersion(deployments, input)
}

func (h Handler) In(input *models.Inputs, outDir string) (*models.Output, error) {
	deploy, ok := input.Version.(gitlab.Deployment)
	if !ok {
		return nil, fmt.Errorf("failed to cast version %v to deployment", input.Version)
	}

	versionJSON, err := json.Marshal(input.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize version to JSON: %w", err)
	}

	err = os.WriteFile(outDir+"/metadata.json", versionJSON, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to wrte version to output dir %q: %w", outDir, err)
	}

	// Get the full user payload to get the user email
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: input.Source.Auth.Token,
		Url:   input.Source.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	user, err := internal.GetUser(deploy.User, client)
	if err != nil {
		return nil, err
	}

	return &models.Output{
		Version: input.Version,
		Metadata: models.Metadatas{
			{Name: "id", Value: strconv.FormatInt(int64(deploy.ID), 10)},
			{Name: "iid", Value: strconv.FormatInt(int64(deploy.IID), 10)},
			{Name: "status", Value: deploy.Status},
			{Name: "ref", Value: deploy.Ref},
			{Name: "sha", Value: deploy.SHA},
			{Name: "deployable_name", Value: deploy.Deployable.Name},
			{Name: "deployable_ref", Value: deploy.Deployable.Ref},
			{Name: "deployable_stage", Value: deploy.Deployable.Stage},
			{Name: "deployable_started_at", Value: deploy.Deployable.StartedAt.String()},
			{Name: "deployable_status", Value: deploy.Deployable.Status},
			{Name: "deployable_tag", Value: strconv.FormatBool(deploy.Deployable.Tag)},
			{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(deploy.Deployable.Pipeline.ID), 10)},
			{Name: "deployable_pipeline_ref", Value: deploy.Deployable.Pipeline.Ref},
			{Name: "deployable_pipeline_sha", Value: deploy.Deployable.Pipeline.SHA},
			{Name: "deployable_pipeline_status", Value: deploy.Deployable.Pipeline.Status},
			{Name: "deployable_pipeline_updated_at", Value: deploy.Deployable.Pipeline.UpdatedAt.String()},
			{Name: "deployable_commit_author_name", Value: deploy.Deployable.Commit.AuthorName},
			{Name: "deployable_commit_author_email", Value: deploy.Deployable.Commit.AuthorEmail},
			{Name: "deployable_commit_title", Value: deploy.Deployable.Commit.Title},
			{Name: "deployable_commit_message", Value: deploy.Deployable.Commit.Message},
			{Name: "deployable_commit_short_id", Value: deploy.Deployable.Commit.ShortID},
			{Name: "deployable_commit_created_at", Value: deploy.Deployable.Commit.CreatedAt.String()},
			{Name: "deployable_commit_updated_at", Value: deploy.Deployable.Commit.LastPipeline.UpdatedAt.String()},
			{Name: "environment_id", Value: strconv.FormatInt(int64(deploy.Environment.ID), 10)},
			{Name: "environment_name", Value: deploy.Environment.Name},
			{Name: "environment_external_url", Value: deploy.Environment.ExternalURL},
			{Name: "user_username", Value: user.Username},
			{Name: "user_name", Value: user.Name},
			{Name: "user_email", Value: user.Email},
			{Name: "user_public_email", Value: user.PublicEmail},
		},
	}, nil
}

func (h Handler) Out(input *models.Inputs, outDir string) (*models.Output, error) {
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: input.Source.Auth.Token,
		Url:   input.Source.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	switch input.Params.Deployments.Action {
	case "create":
		if input.Params.Deployments.Create == nil {
			return nil, fmt.Errorf("missing params.create source")
		}

		deploy, _, err := client.Deployments.CreateProjectDeployment(
			input.Source.ProjectID, &gitlab.CreateProjectDeploymentOptions{
				Environment: input.Source.Deployments.Environment,
				SHA:         &input.Params.Deployments.Create.SHA,
				Ref:         &input.Params.Deployments.Create.Ref,
				Tag:         &input.Params.Deployments.Create.Tag,
				Status:      (*gitlab.DeploymentStatusValue)(&input.Params.Deployments.Create.Status),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create deployment: %w", err)
		}

		deployJSON, err := json.MarshalIndent(deploy, "", "  ")
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			fmt.Fprintf(os.Stderr, "failed to serialize deployment payload to JSON: %s\n", err.Error())
		}

		err = os.WriteFile(outDir+"/metadata.json", deployJSON, 0666)
		if err != nil {
			// If we fail to write the JSON, we should still try to send the metadata JSON
			fmt.Fprintf(os.Stderr, "failed to write metadata to %q: %s\n", outDir+"/metadata.json", err.Error())
		}

		return &models.Output{
			Version: deploy,
			Metadata: models.Metadatas{
				{Name: "id", Value: strconv.FormatInt(int64(deploy.ID), 10)},
				{Name: "iid", Value: strconv.FormatInt(int64(deploy.IID), 10)},
				{Name: "status", Value: deploy.Status},
				{Name: "ref", Value: deploy.Ref},
				{Name: "sha", Value: deploy.SHA},
				{Name: "deployable_name", Value: deploy.Deployable.Name},
				{Name: "deployable_ref", Value: deploy.Deployable.Ref},
				{Name: "deployable_stage", Value: deploy.Deployable.Stage},
				{Name: "deployable_started_at", Value: deploy.Deployable.StartedAt.String()},
				{Name: "deployable_status", Value: deploy.Deployable.Status},
				{Name: "deployable_tag", Value: strconv.FormatBool(deploy.Deployable.Tag)},
				{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(deploy.Deployable.Pipeline.ID), 10)},
				{Name: "deployable_pipeline_ref", Value: deploy.Deployable.Pipeline.Ref},
				{Name: "deployable_pipeline_sha", Value: deploy.Deployable.Pipeline.SHA},
				{Name: "deployable_pipeline_status", Value: deploy.Deployable.Pipeline.Status},
				{Name: "deployable_pipeline_updated_at", Value: deploy.Deployable.Pipeline.UpdatedAt.String()},
				{Name: "deployable_commit_author_name", Value: deploy.Deployable.Commit.AuthorName},
				{Name: "deployable_commit_author_email", Value: deploy.Deployable.Commit.AuthorEmail},
				{Name: "deployable_commit_title", Value: deploy.Deployable.Commit.Title},
				{Name: "deployable_commit_message", Value: deploy.Deployable.Commit.Message},
				{Name: "deployable_commit_short_id", Value: deploy.Deployable.Commit.ShortID},
				{Name: "deployable_commit_created_at", Value: deploy.Deployable.Commit.CreatedAt.String()},
				{Name: "deployable_commit_updated_at", Value: deploy.Deployable.Commit.LastPipeline.UpdatedAt.String()},
				{Name: "environment_id", Value: strconv.FormatInt(int64(deploy.Environment.ID), 10)},
				{Name: "environment_name", Value: deploy.Environment.Name},
				{Name: "environment_external_url", Value: deploy.Environment.ExternalURL},
			},
		}, nil

	case "update":
		if input.Params.Deployments.Update == nil {
			return nil, fmt.Errorf("missing params.update source")
		}

		deploy, ok := input.Version.(gitlab.Deployment)
		if !ok {
			return nil, fmt.Errorf("failed to cast version %v to deployment", input.Version)
		}

		updatedDeploy, _, err := client.Deployments.UpdateProjectDeployment(
			input.Source.ProjectID, deploy.ID, &gitlab.UpdateProjectDeploymentOptions{
				Status: (*gitlab.DeploymentStatusValue)(&input.Params.Deployments.Update.Status),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update deployment with id %q: %w", deploy.ID, err)
		}

		user, err := internal.GetUser(deploy.User, client)
		if err != nil {
			return nil, err
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

		return &models.Output{
			Version: updatedDeploy,
			Metadata: models.Metadatas{
				{Name: "id", Value: strconv.FormatInt(int64(deploy.ID), 10)},
				{Name: "iid", Value: strconv.FormatInt(int64(deploy.IID), 10)},
				{Name: "status", Value: deploy.Status},
				{Name: "ref", Value: deploy.Ref},
				{Name: "sha", Value: deploy.SHA},
				{Name: "deployable_name", Value: deploy.Deployable.Name},
				{Name: "deployable_ref", Value: deploy.Deployable.Ref},
				{Name: "deployable_stage", Value: deploy.Deployable.Stage},
				{Name: "deployable_started_at", Value: deploy.Deployable.StartedAt.String()},
				{Name: "deployable_status", Value: deploy.Deployable.Status},
				{Name: "deployable_tag", Value: strconv.FormatBool(deploy.Deployable.Tag)},
				{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(deploy.Deployable.Pipeline.ID), 10)},
				{Name: "deployable_pipeline_ref", Value: deploy.Deployable.Pipeline.Ref},
				{Name: "deployable_pipeline_sha", Value: deploy.Deployable.Pipeline.SHA},
				{Name: "deployable_pipeline_status", Value: deploy.Deployable.Pipeline.Status},
				{Name: "deployable_pipeline_updated_at", Value: deploy.Deployable.Pipeline.UpdatedAt.String()},
				{Name: "deployable_commit_author_name", Value: deploy.Deployable.Commit.AuthorName},
				{Name: "deployable_commit_author_email", Value: deploy.Deployable.Commit.AuthorEmail},
				{Name: "deployable_commit_title", Value: deploy.Deployable.Commit.Title},
				{Name: "deployable_commit_message", Value: deploy.Deployable.Commit.Message},
				{Name: "deployable_commit_short_id", Value: deploy.Deployable.Commit.ShortID},
				{Name: "deployable_commit_created_at", Value: deploy.Deployable.Commit.CreatedAt.String()},
				{Name: "deployable_commit_updated_at", Value: deploy.Deployable.Commit.LastPipeline.UpdatedAt.String()},
				{Name: "environment_id", Value: strconv.FormatInt(int64(deploy.Environment.ID), 10)},
				{Name: "environment_name", Value: deploy.Environment.Name},
				{Name: "environment_external_url", Value: deploy.Environment.ExternalURL},
				{Name: "user_username", Value: user.Username},
				{Name: "user_name", Value: user.Name},
				{Name: "user_email", Value: user.Email},
				{Name: "user_public_email", Value: user.PublicEmail},
			},
		}, nil

	case "delete":
		deploy, ok := input.Version.(gitlab.Deployment)
		if !ok {
			return nil, fmt.Errorf("failed to cast version %v to deployment", input.Version)
		}

		_, err := client.Deployments.DeleteProjectDeployment(input.Source.ProjectID, deploy.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to update deployment with id %q: %w", deploy.ID, err)
		}

		return &models.Output{
			Version:  nil,
			Metadata: models.Metadatas{},
		}, nil

	default:
		return nil, fmt.Errorf("invalid params.action parameter, accepted values are: create, update, delete")
	}
}
