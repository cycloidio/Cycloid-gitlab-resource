package deployments

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h Handler) Check() error {
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	var options = &gitlab.ListProjectDeploymentsOptions{
		Sort:        gitlab.Ptr("desc"), // version should be oldest first
		Environment: h.cfg.Source.Environment,
		Status:      h.cfg.Source.Status,
	}

	deployments, _, err := client.Deployments.ListProjectDeployments(h.cfg.Source.ProjectID, options, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch deployments from gitlab API: %w", err)
	}

	if len(deployments) == 0 {
		_, err := fmt.Fprintf(h.stdout, "[]")
		if err != nil {
			return fmt.Errorf("failed to output to h.cfg.stdout: %w", err)
		}

		return nil
	}

	if h.cfg.Source.Mode == "yield" {
		// Always return the oldest deployment with a timestamp
		version := DeploymentToVersion(deployments[0])
		version["check_timestamp"] = time.Now().String()
		return OutputJSON(h.stdout, version)
	}

	var versions []map[string]string
	if h.cfg.Version == nil {
		versions = []map[string]string{DeploymentToVersion(deployments[0])}
	} else {

		currentIDStr, ok := h.cfg.Version["id"]
		if !ok {
			return fmt.Errorf("failed to get current id from version %v", h.cfg.Version)
		}

		currentID, err := strconv.ParseInt(currentIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse version id %q: %w", currentIDStr, err)
		}

		if i := slices.IndexFunc(deployments, func(d *gitlab.Deployment) bool {
			return d.ID == int(currentID)
		}); i != -1 {
			versions = DeploymentsToVersion(deployments[:i])
		} else {
			versions = DeploymentsToVersion(deployments)
		}
	}

	return OutputJSON(h.stdout, versions)
}
