package deployments

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
)

func (h Handler) In(outDir string) error {
	versionJSON, err := json.Marshal(h.cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to serialize version to JSON: %w", err)
	}

	err = os.WriteFile(outDir+"/metadata.json", versionJSON, 0666)
	if err != nil {
		return fmt.Errorf("failed to write version to output dir %q: %w", outDir, err)
	}

	// Get the full user payload to get the user email
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	user, err := internal.GetUser(h.cfg.Version.User, client)
	if err != nil {
		return err
	}

	output := &models.Output{
		Version: h.cfg.Version,
		Metadata: models.Metadatas{
			{Name: "id", Value: strconv.FormatInt(int64(h.cfg.Version.ID), 10)},
			{Name: "iid", Value: strconv.FormatInt(int64(h.cfg.Version.IID), 10)},
			{Name: "status", Value: h.cfg.Version.Status},
			{Name: "ref", Value: h.cfg.Version.Ref},
			{Name: "sha", Value: h.cfg.Version.SHA},
			{Name: "deployable_name", Value: h.cfg.Version.Deployable.Name},
			{Name: "deployable_ref", Value: h.cfg.Version.Deployable.Ref},
			{Name: "deployable_stage", Value: h.cfg.Version.Deployable.Stage},
			{Name: "deployable_started_at", Value: h.cfg.Version.Deployable.StartedAt.String()},
			{Name: "deployable_status", Value: h.cfg.Version.Deployable.Status},
			{Name: "deployable_tag", Value: strconv.FormatBool(h.cfg.Version.Deployable.Tag)},
			{Name: "deployable_pipeline_id", Value: strconv.FormatInt(int64(h.cfg.Version.Deployable.Pipeline.ID), 10)},
			{Name: "deployable_pipeline_ref", Value: h.cfg.Version.Deployable.Pipeline.Ref},
			{Name: "deployable_pipeline_sha", Value: h.cfg.Version.Deployable.Pipeline.SHA},
			{Name: "deployable_pipeline_status", Value: h.cfg.Version.Deployable.Pipeline.Status},
			{Name: "deployable_pipeline_updated_at", Value: h.cfg.Version.Deployable.Pipeline.UpdatedAt.String()},
			{Name: "deployable_commit_author_name", Value: h.cfg.Version.Deployable.Commit.AuthorName},
			{Name: "deployable_commit_author_email", Value: h.cfg.Version.Deployable.Commit.AuthorEmail},
			{Name: "deployable_commit_title", Value: h.cfg.Version.Deployable.Commit.Title},
			{Name: "deployable_commit_message", Value: h.cfg.Version.Deployable.Commit.Message},
			{Name: "deployable_commit_short_id", Value: h.cfg.Version.Deployable.Commit.ShortID},
			{Name: "deployable_commit_created_at", Value: h.cfg.Version.Deployable.Commit.CreatedAt.String()},
			{Name: "deployable_commit_updated_at", Value: h.cfg.Version.Deployable.Commit.LastPipeline.UpdatedAt.String()},
			{Name: "environment_id", Value: strconv.FormatInt(int64(h.cfg.Version.Environment.ID), 10)},
			{Name: "environment_name", Value: h.cfg.Version.Environment.Name},
			{Name: "environment_external_url", Value: h.cfg.Version.Environment.ExternalURL},
			{Name: "user_username", Value: user.Username},
			{Name: "user_name", Value: user.Name},
			{Name: "user_email", Value: user.Email},
			{Name: "user_public_email", Value: user.PublicEmail},
		},
	}

	return OutputJSON(h.stdout, output)
}
