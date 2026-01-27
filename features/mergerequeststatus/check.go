package mergerequeststatus

import (
	"fmt"
	"slices"

	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) Check() error {
	var options = &gitlab.GetMergeRequestsOptions{}
	mr, _, err := h.glab.MergeRequests.GetMergeRequest(
		h.cfg.Source.ProjectID,
		h.cfg.Source.MergeRequestStatus.MergeRequestIID,
		options,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get mr with iid %d and project id %q from API: %w",
			h.cfg.Source.MergeRequestStatus.MergeRequestIID,
			h.cfg.Source.ProjectID,
			err,
		)
	}

	if slices.Contains(h.cfg.Source.MergeRequestStatus.State, models.MergeRequestState(mr.State)) || len(h.cfg.Source.MergeRequestStatus.State) == 0 {
		return internal.OutputJSON(h.stdout, []map[string]string{MergeRequestToVersion(mr)})
	}

	return internal.OutputJSON(h.stdout, []map[string]string{})
}
