package deployments

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func OutputJSON(stdout io.Writer, output any) error {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize output to JSON: %w", err)
	}

	_, err = stdout.Write(data)
	if err != nil {
		return fmt.Errorf("failed to output to stdout: %w", err)
	}

	return nil
}

func DeploymentToVersion(deploy *gitlab.Deployment) map[string]string {
	var version = make(map[string]string)
	version["id"] = strconv.FormatInt(int64(deploy.ID), 10)
	version["iid"] = strconv.FormatInt(int64(deploy.IID), 10)
	version["status"] = deploy.Status
	version["ref"] = deploy.Ref
	version["sha"] = deploy.SHA
	version["sha"] = deploy.SHA
	version["user_id"] = strconv.FormatInt(int64(deploy.User.ID), 10)
	version["user_username"] = deploy.User.Username
	version["user_name"] = deploy.User.Name
	version["user_web_url"] = deploy.User.WebURL
	version["environment_id"] = strconv.FormatInt(int64(deploy.Environment.ID), 10)
	version["environment_name"] = deploy.Environment.Name
	version["environment_tier"] = deploy.Environment.Tier
	version["environment_external_url"] = deploy.Environment.ExternalURL
	return version
}

func DeploymentsToVersion(deployments []*gitlab.Deployment) []map[string]string {
	var version = make([]map[string]string, len(deployments))
	for i, deploy := range deployments {
		version[i] = DeploymentToVersion(deploy)
	}

	return version
}
