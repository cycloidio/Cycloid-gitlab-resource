package mergerequeststatus

import (
	"fmt"

	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) In(outDir string) error {
	var options = &gitlab.GetMergeRequestsOptions{}
	mr, _, err := h.glab.MergeRequests.GetMergeRequest(
		h.cfg.Source.ProjectID,
		h.cfg.Source.MergeRequestStatus.MergeRequestID,
		options,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get mr with iid %q and project id %q from API: %w",
			h.cfg.Source.MergeRequestStatus.MergeRequestID,
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
