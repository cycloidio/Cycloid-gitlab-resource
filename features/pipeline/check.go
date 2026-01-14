package pipeline

import (
	"fmt"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) Check() error {
	glab, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	if h.cfg.Source.Pipeline.CheckMode == nil {
		h.cfg.Source.Pipeline.CheckMode = gitlab.Ptr(models.Get)
	}

	switch *h.cfg.Source.Pipeline.CheckMode {
	case models.Get:
		pipeline, _, err := glab.Pipelines.GetLatestPipeline(h.cfg.Source.ProjectID,
			&gitlab.GetLatestPipelineOptions{Ref: h.cfg.Source.Pipeline.Ref},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to fetch latest pipeline in project id %q with ref %v: %w",
				h.cfg.Source.ProjectID, h.cfg.Source.Pipeline.Ref,
				err,
			)
		}

		return internal.PrintJSON(h.stdout, PipelinetoVersion(pipeline))
	default:
		return fmt.Errorf("unsuported check mode %q", *h.cfg.Source.Pipeline.CheckMode)
	}
}
