package pipeline

import (
	"fmt"
	"path"
	"slices"
	"strconv"

	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	"github.com/sanity-io/litter"
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
			return fmt.Errorf("either params.pipeline.merge_request_iid or params.pipeline.ref is required to trigger a pipeline")
		}

		h.logger.Info("create pipeline succeeded", "id", pipeline.ID, "url", pipeline.WebURL)

		pipelineStatus := pipeline.Status
		for !slices.Contains([]string{"success", "failed", "canceled", "skipped", "manual"}, pipelineStatus) {
			// watch and refresh pipeline status
			pipeline, _, err := h.glab.Pipelines.GetPipeline(h.cfg.Source.ProjectID, pipeline.ID)
			if err != nil {
				return fmt.Errorf("failed to get pipeline status for %q: %w", strconv.Itoa(pipeline.ID), err)
			}
			pipelineStatus = pipeline.Status
			h.logger.Debug("pipeline", "status", pipelineStatus)

			// list the pending jobs of the current pipeline
			jobsPending, _, err := h.glab.Jobs.ListPipelineJobs(h.cfg.Source.ProjectID, pipeline.ID,
				nil,
				// &gitlab.ListJobsOptions{
				// 	ListOptions: gitlab.ListOptions{},
				// 	Scope: &[]gitlab.BuildStateValue{
				// 		gitlab.Running,
				// 		gitlab.Pending,
				// 		gitlab.Preparing,
				// 		gitlab.Created,
				// 	},
				// },
			)
			if err != nil {
				return fmt.Errorf("failed to list jobs from pipeline %q: %w", strconv.Itoa(pipeline.ID), err)
			}

			// List jobs that trigger child pipelines
			jobsTriggerBridges, _, err := h.glab.Jobs.ListPipelineBridges(h.cfg.Source.ProjectID, pipeline.ID, &gitlab.ListJobsOptions{})
			if err != nil {
				return fmt.Errorf("failed to list trigger jobs from pipeline %q: %w", strconv.Itoa(pipeline.ID), err)
			}

			for _, bridge := range jobsTriggerBridges {
				litter.D(bridge)
				if bridge.DownstreamPipeline == nil {
					h.logger.Debug("trigger job not started yet", "status", bridge.Status, "downstream pipeline", bridge.DownstreamPipeline)
					continue
				}

				childJobs, _, err := h.glab.Jobs.ListPipelineJobs(bridge.DownstreamPipeline.ProjectID, bridge.DownstreamPipeline.ID, &gitlab.ListJobsOptions{})
				if err != nil {
					return fmt.Errorf("failed to list jobs from child pipeline: %w", err)
				}
				litter.D(childJobs)

				jobsPending = append(jobsPending, childJobs...)
			}

			var jobsRunning = []*gitlab.Job{}
			for _, job := range jobsPending {
				if job.Status == "running" {
					jobsRunning = append(jobsRunning, job)
				}
			}

			jobNames := make([]string, len(jobsPending))
			for i, job := range jobsPending {
				jobNames[i] = job.Name
			}
			h.logger.Debug("Jobs pending", "jobs", jobNames)

			for _, job := range jobsRunning {
				err := h.traceJob(job)
				if err != nil {
					return fmt.Errorf("failed to fetch logs for job %q: %w", job.Name, err)
				}
			}
		}

		// refresh the status
		pipeline, _, err = h.glab.Pipelines.GetPipeline(h.cfg.Source.ProjectID, pipeline.ID)
		if err != nil {
			return fmt.Errorf("failed to get pipeline status for %q: %w", strconv.Itoa(pipeline.ID), err)
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
