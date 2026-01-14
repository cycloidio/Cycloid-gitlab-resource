package deployments

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
	metadataPath := path.Join(outDir, "metadata.json")
	versionJSON, err := json.Marshal(h.cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to serialize version to JSON: %w", err)
	}

	err = os.WriteFile(metadataPath, versionJSON, 0666)
	if err != nil {
		return fmt.Errorf("failed to write version to output dir %q: %w", outDir, err)
	}

	deployID, err := strconv.ParseInt(h.cfg.Version["id"], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to cast current deployment id %q: %w", h.cfg.Version["id"], err)
	}

	deploy, _, err := h.glab.Deployments.GetProjectDeployment(h.cfg.Source.ProjectID, int(deployID))
	if err != nil {
		return fmt.Errorf("failed to fetch deployment with id %q: %w", h.cfg.Version["id"], err)
	}

	if h.cfg.Source.Status != nil {
		if deploy.Status != *h.cfg.Source.Status {
			return fmt.Errorf("deployment with id %d has status %q where status %q is requested", deploy.ID, deploy.Status, *h.cfg.Source.Status)
		}
	}

	newVersion := DeploymentToVersion(deploy)
	h.cfg.Version = newVersion
	var metdatas = models.Metadatas{
		{Name: "id", Value: h.cfg.Version["id"]},
		{Name: "status", Value: h.cfg.Version["status"]},
		{Name: "ref", Value: h.cfg.Version["ref"]},
		{Name: "sha", Value: h.cfg.Version["sha"]},
		{Name: "environment_id", Value: h.cfg.Version["environment_id"]},
		{Name: "environment_name", Value: h.cfg.Version["environment_name"]},
		{Name: "environment_external_url", Value: h.cfg.Version["environment_tier"]},
	}

	if userIDStr, ok := h.cfg.Version["user_id"]; ok {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse user id %q: %w", userIDStr, err)
		}

		user, err := internal.GetUser(int(userID), h.glab)
		if err != nil {
			return err
		}

		metdatas = append(metdatas,
			models.Metadatas{
				{Name: "user_username", Value: user.Username},
				{Name: "user_email", Value: user.Email},
			}...,
		)
	}

	output := &models.Output{
		Version:  h.cfg.Version,
		Metadata: metdatas,
	}

	return internal.OutputJSON(h.stdout, output)
}
