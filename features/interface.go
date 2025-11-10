package features

import (
	"fmt"
	"strings"

	"github.com/cycloidio/gitlab-resource/features/deployments"
	"github.com/cycloidio/gitlab-resource/models"
)

var (
	AvailableFeatures = [...]string{
		"deployments",
	}
)

type Feature interface {
	Validate(input *models.Inputs) error
	Check(input *models.Inputs) (any, error)
	In(input *models.Inputs, outDir string) (*models.Output, error)
	Out(input *models.Inputs, outDir string) (*models.Output, error)
}

func NewFeatureHandler(feature string) (Feature, error) {
	switch strings.ToLower(feature) {
	case "deployments":
		return deployments.NewHandler(), nil
	default:
		return nil, fmt.Errorf("feature %q does not exists, available ones are: %v", feature, AvailableFeatures)
	}
}
