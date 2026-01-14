package pipeline

import (
	"fmt"
	"io"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (h *Handler) traceJob(job *gitlab.Job) error {
	if job == nil {
		return fmt.Errorf("missing job for tracing, job is nil")
	}

	var offset int64
	_, _ = fmt.Fprintf(h.stdout, "=== Logs for job %q ===\n", job.Name)

	for range time.NewTicker(time.Second * 3).C {
		jobStatus, _, err := h.glab.Jobs.GetJob(h.cfg.Source.ProjectID, job.ID)
		if err != nil {
			return fmt.Errorf("failed to query job %q: %w", job.Name, err)
		}

		trace, _, err := h.glab.Jobs.GetTraceFile(h.cfg.Source.ProjectID, job.ID)
		if err != nil {
			return fmt.Errorf("failed to get job logs: %w", err)
		}

		_, _ = io.CopyN(io.Discard, trace, offset)
		traceLen, err := io.Copy(h.stdout, trace)
		if err != nil {
			return fmt.Errorf("failed to write trace to stdout: %w", err)
		}
		offset = offset + traceLen

		switch jobStatus.Status {
		case string(gitlab.Running):
			continue
		case string(gitlab.Pending):
			_, _ = fmt.Fprintf(h.stdout, "%q is pending... waiting for job to start.\n", job.Name)
			continue
		case string(gitlab.Manual):
			_, _ = fmt.Fprintf(h.stdout, "%q is a manual job, skipping", job.Name)
			return nil
		case string(gitlab.Skipped):
			_, _ = fmt.Fprintf(h.stdout, "%q has been skipped.\n", job.Name)
			return nil
		case string(gitlab.Success):
			_, _ = fmt.Fprintf(h.stdout, "%q succeeded.\n", job.Name)
			return nil
		case string(gitlab.Failed):
			_, _ = fmt.Fprintf(h.stdout, "%q job failed.\n", job.Name)
			return nil
		case string(gitlab.Canceled):
			_, _ = fmt.Fprintf(h.stdout, "%q has been cancelled.\n", job.Name)
			return nil
		default:
			return fmt.Errorf("unexpected job status %q from api for job %q", job.Status, job.Name)
		}
	}

	return nil
}
