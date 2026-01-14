package pipeline

import (
	"strconv"

	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func PipelinetoVersion(pipeline *gitlab.Pipeline) map[string]string {
	var version = make(map[string]string)
	version["id"] = strconv.Itoa(pipeline.ID)
	version["iid"] = strconv.Itoa(pipeline.IID)
	version["project_id"] = strconv.Itoa(pipeline.ProjectID)
	version["name"] = pipeline.Name
	version["ref"] = pipeline.Ref
	version["status"] = pipeline.Status
	version["web_url"] = pipeline.WebURL
	version["before_sha"] = pipeline.BeforeSHA
	version["tag"] = strconv.FormatBool(pipeline.Tag)
	if pipeline.CreatedAt != nil {
		version["created_at"] = pipeline.CreatedAt.String()
	}
	if pipeline.UpdatedAt != nil {
		version["updated_at"] = pipeline.UpdatedAt.String()
	}
	if pipeline.StartedAt != nil {
		version["started_at"] = pipeline.StartedAt.String()
	}
	return version
}

func PipelinetoMetadatas(pipeline *gitlab.Pipeline) models.Metadatas {
	metadatas := models.Metadatas{
		{Name: "id", Value: strconv.Itoa(pipeline.ID)},
		{Name: "iid", Value: strconv.Itoa(pipeline.IID)},
		{Name: "project_id", Value: strconv.Itoa(pipeline.ProjectID)},
		{Name: "name", Value: pipeline.Name},
		{Name: "ref", Value: pipeline.Ref},
		{Name: "status", Value: pipeline.Status},
		{Name: "web_url", Value: pipeline.WebURL},
		{Name: "before_sha", Value: pipeline.BeforeSHA},
		{Name: "tag", Value: strconv.FormatBool(pipeline.Tag)},
	}
	if pipeline.CreatedAt != nil {
		metadatas = append(metadatas, models.Metadata{Name: "created_at", Value: pipeline.CreatedAt.String()})
	}
	if pipeline.UpdatedAt != nil {
		metadatas = append(metadatas, models.Metadata{Name: "updated_at", Value: pipeline.UpdatedAt.String()})
	}
	if pipeline.StartedAt != nil {
		metadatas = append(metadatas, models.Metadata{Name: "started_at", Value: pipeline.StartedAt.String()})
	}
	return metadatas
}
