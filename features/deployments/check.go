package deployments

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/cycloidio/gitlab-resource/internal"
	gitlab "gitlab.com/gitlab-org/api/h.glab-go"
)

func (h Handler) Check() error {
	var options = &gitlab.ListProjectDeploymentsOptions{
		Sort:        gitlab.Ptr("desc"), // version should be oldest first
		Environment: h.cfg.Source.Environment,
		Status:      h.cfg.Source.Status,
	}

	deployments, _, err := h.glab.Deployments.ListProjectDeployments(h.cfg.Source.ProjectID, options, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch deployments from gitlab API: %w", err)
	}
	deploymentsLen := len(deployments)

	if deploymentsLen == 0 {
		_, err := fmt.Fprintf(h.stdout, "[]")
		if err != nil {
			return fmt.Errorf("failed to output to h.cfg.stdout: %w", err)
		}

		return nil
	}

	var versions = []map[string]string{}
	if h.cfg.Source.Mode == "yield" {
		version := DeploymentToVersion(deployments[deploymentsLen-1])
		versions = append(versions, version)
	} else if h.cfg.Version == nil {
		versions = []map[string]string{DeploymentToVersion(deployments[deploymentsLen-1])}
	} else {
		if h.cfg.Version["status"] != "running" {
			// don't return new version if the current one doesn't match the status
			return internal.OutputJSON(h.stdout, versions)
		}

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

	return internal.OutputJSON(h.stdout, versions)
}
