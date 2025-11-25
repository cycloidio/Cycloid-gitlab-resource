package environments

import (
	"fmt"
	"slices"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h Handler) Check() error {
	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	var options = &gitlab.ListEnvironmentsOptions{
		Name:   h.cfg.Source.Environment.Name,
		Search: h.cfg.Source.Environment.Search,
		States: h.cfg.Source.Environment.States,
	}

	envs, _, err := glab.Environments.ListEnvironments(h.cfg.Source.ProjectID, options)
	if err != nil {
		return fmt.Errorf("failed to fetch environments from gitlab API: %w", err)
	}
	envLen := len(envs)

	if envLen == 0 {
		_, err := h.stdout.Write([]byte("[]"))
		if err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}

		return nil
	}

	var versions = []map[string]string{}
	if h.cfg.Version == nil {
		versions = []map[string]string{EnvironmentToVersion(envs[0])}
	} else {
		if h.cfg.Source.Environment.States != nil {
			if h.cfg.Version["state"] != *h.cfg.Source.Environment.States {
				// don't return new version if the current one doesn't match the status
				// this can occur after a put step
				return internal.OutputJSON(h.stdout, versions)
			}
		}

		currentIDStr, ok := h.cfg.Version["id"]
		if !ok {
			return fmt.Errorf("failed to get current id from version %v", h.cfg.Version)
		}

		currentID, err := strconv.ParseInt(currentIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse version id %q, %w", currentIDStr, err)
		}

		if i := slices.IndexFunc(envs, func(e *gitlab.Environment) bool {
			return e.ID == int(currentID)
		}); i != -1 {
			versions = EnvironmentsToVersion(envs[:i])
		} else {
			versions = EnvironmentsToVersion(envs)
		}
	}

	return internal.OutputJSON(h.stdout, versions)
}
