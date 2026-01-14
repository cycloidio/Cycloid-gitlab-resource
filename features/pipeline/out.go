package pipeline

import (
	"fmt"
	"path"
	"strconv"

	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) Out(outDir string) error {
	metadataPath := path.Join(outDir, "metadata.json")
	_ = metadataPath

	switch h.cfg.Params.Action {
	case models.ActionCreate:
		pipeline, _, err := h.glab.Pipelines.CreatePipeline(
			h.cfg.Source.ProjectID, &gitlab.CreatePipelineOptions{
				Ref:       &h.cfg.Params.Ref,
				Variables: h.cfg.Params.Variables,
				Inputs:    h.cfg.Params.Inputs,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to trigger pipeline with ref %q: %w", h.cfg.Params.Ref, err)
		}

		h.logger.Info("Created pipeline succeeded", "id", pipeline.ID, "url", pipeline.WebURL)

		for {
			jobsPending, _, err := h.glab.Jobs.ListPipelineJobs(h.cfg.Source.ProjectID, pipeline.ID, &gitlab.ListJobsOptions{
				Scope: &[]gitlab.BuildStateValue{
					gitlab.Running,
					gitlab.Pending,
					gitlab.Preparing,
					gitlab.Created,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list jobs from pipeline %q: %w", strconv.Itoa(pipeline.ID), err)
			}

			jobsRunning, _, err := h.glab.Jobs.ListPipelineJobs(h.cfg.Source.ProjectID, pipeline.ID, &gitlab.ListJobsOptions{
				Scope: &[]gitlab.BuildStateValue{
					gitlab.Running,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list jobs from pipeline %q: %w", strconv.Itoa(pipeline.ID), err)
			}

			jobNames := make([]string, len(jobsPending))
			for i, job := range jobsPending {
				jobNames[i] = job.Name
			}
			h.logger.Debug("Jobs pending", "jobs", jobNames)

			if len(jobsPending) == 0 {
				h.logger.Debug("jobs are done", "jobs", jobsPending)
				break
			}

			for _, job := range jobsRunning {
				err := h.traceJob(job)
				if err != nil {
					return fmt.Errorf("failed to fetch logs for job %q: %w", job.Name, err)
				}
			}
		}

		output := &models.Output{
			Version:  PipelinetoVersion(pipeline),
			Metadata: PipelinetoMetadatas(pipeline),
		}

		return internal.OutputJSON(h.stdout, output)
	default:
		return fmt.Errorf("invalid action %q, available ones are: %q", h.cfg.Params.Action, models.ActionCreate)
	}
}
