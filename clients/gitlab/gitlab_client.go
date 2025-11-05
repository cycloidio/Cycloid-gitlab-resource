package gitlabclient

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabConfig struct {
	Token string
	Url   string
}

func NewGitlabClient(config *GitlabConfig) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.Url))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}

	return client, nil
}
