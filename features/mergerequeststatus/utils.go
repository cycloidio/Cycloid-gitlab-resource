package mergerequeststatus

import (
	"strconv"

	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func MergeRequestToVersion(mr *gitlab.MergeRequest) map[string]string {
	var version = make(map[string]string)
	version["iid"] = strconv.FormatInt(int64(mr.IID), 10)
	version["state"] = mr.State
	version["author_username"] = mr.Author.Username
	version["author_id"] = strconv.FormatInt(int64(mr.Author.ID), 10)
	version["project_id"] = strconv.FormatInt(int64(mr.ProjectID), 10)
	version["title"] = mr.Title
	version["description"] = mr.Description
	return version
}

func MergeRequestToMetadatas(mr *gitlab.MergeRequest) models.Metadatas {
	return models.Metadatas{
		{Name: "iid", Value: strconv.Itoa(mr.IID)},
		{Name: "state", Value: mr.State},
		{Name: "author_username", Value: mr.Author.Username},
		{Name: "author_id", Value: strconv.Itoa(mr.Author.ID)},
		{Name: "project_id", Value: strconv.Itoa(mr.ProjectID)},
		{Name: "title", Value: mr.Title},
		{Name: "description", Value: mr.Description},
	}
}
