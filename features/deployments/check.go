package deployments

import (
	"fmt"
	"slices"

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

	if h.cfg.Version == nil {
		return OutputJSON(h.stdout, deployments[0])
	} else {
		var versions []*gitlab.Deployment
		if i := slices.IndexFunc(deployments, func(d *gitlab.Deployment) bool {
			return d.ID == h.cfg.Version.ID
		}); i != -1 {
			versions = deployments[:i]
		} else {
			versions = deployments
		}

		return OutputJSON(h.stdout, versions)
	}
}
