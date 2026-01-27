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
		pipelineInput := make(gitlab.PipelineInputsOption)
		for key, value := range h.cfg.Params.Inputs {
			switch v := value.(type) {
			case string:
				pipelineInput[key] = gitlab.NewPipelineInputValue(v)
			case []string:
				pipelineInput[key] = gitlab.NewPipelineInputValue(v)
			case float64:
				pipelineInput[key] = gitlab.NewPipelineInputValue(v)
			case int:
				pipelineInput[key] = gitlab.NewPipelineInputValue(v)
			case bool:
				pipelineInput[key] = gitlab.NewPipelineInputValue(v)
			default:
				return fmt.Errorf("invalid type %T for input with key: %q, only string, []string, bool, int and float accepted by gitlab api", v, key)
			}
		}

		var pipeline *gitlab.Pipeline
		var err error
		if h.cfg.Params.MergeRequestIID != nil {
			h.logger.Debug("trigger mr pipeline", "id", h.cfg.Params.MergeRequestIID, "project", h.cfg.Source.ProjectID)
			ppInfo, _, err := h.glab.MergeRequests.CreateMergeRequestPipeline(
				h.cfg.Source.ProjectID, *h.cfg.Params.MergeRequestIID,
			)
			if err != nil {
				return fmt.Errorf("failed to trigger pipeline for merge request %d: %w", *h.cfg.Params.MergeRequestIID, err)
			}

			pipeline, _, err = h.glab.Pipelines.GetPipeline(ppInfo.ProjectID, ppInfo.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch pipeline info for merge request pipeline %q: %w", ppInfo.ID, err)
			}
		} else if h.cfg.Params.Ref != nil {
			pipeline, _, err = h.glab.Pipelines.CreatePipeline(
				h.cfg.Source.ProjectID, &gitlab.CreatePipelineOptions{
					Ref:       h.cfg.Params.Ref,
					Variables: h.cfg.Params.Variables,
					Inputs:    pipelineInput,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to trigger pipeline with ref %q: %w", *h.cfg.Params.Ref, err)
			}
		} else {
			return fmt.Errorf("Either params.pipeline.merge_request or params.pipeline.ref is required to trigger a pipeline.")
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
