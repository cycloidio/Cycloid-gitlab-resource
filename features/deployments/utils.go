package deployments

import (
	"fmt"
	"slices"

	"github.com/cycloidio/gitlab-resource/models"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GetVersion(deployments []*gitlab.Deployment, input *models.Inputs) (any, error) {
	deploymentsLen := len(deployments)
	if deploymentsLen == 0 {
		return []any{}, nil
	}

	if input.Version == nil {
		return deployments[deploymentsLen-1], nil
	}

	deploymentVersion, ok := input.Version.(*gitlab.Deployment)
	if !ok {
		return nil, fmt.Errorf("cannot read current version %v", input.Version)
	}

	if i := slices.IndexFunc(deployments, func(d *gitlab.Deployment) bool {
		return d.ID == deploymentVersion.ID
	}); i != -1 {
		return deployments[i:], nil
	} else {
		return deployments, nil
	}
}
