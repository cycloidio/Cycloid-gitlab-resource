package mergerequeststatus

import (
	"fmt"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) In(outDir string) error {
	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	var options = &gitlab.GetMergeRequestsOptions{}
	mr, _, err := glab.MergeRequests.GetMergeRequest(
		h.cfg.Source.ProjectID,
		h.cfg.Source.MergeRequestStatus.MergeRequestIID,
		options,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get mr with iid %q and project id %q from API: %w",
			h.cfg.Source.MergeRequestStatus.MergeRequestIID,
			h.cfg.Source.ProjectID,
			err,
		)
	}

	output := models.Output{
		Version:  h.cfg.Version,
		Metadata: MergeRequestToMetadatas(mr),
	}

	return internal.OutputJSON(h.stdout, output)
}
