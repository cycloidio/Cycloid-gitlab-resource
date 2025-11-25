package features

import (
	"fmt"
	"io"
	"strings"

	"github.com/cycloidio/gitlab-resource/features/deployments"
	"github.com/cycloidio/gitlab-resource/features/environments"
)

var (
	AvailableFeatures = [...]string{
		"deployments",
		"environments",
	}
)

type Feature interface {
	Check() error
	In(outDir string) error
	Out(outDir string) error
}

func NewFeatureHandler(stdout, stderr io.Writer, feature string, input []byte) (Feature, error) {
	switch strings.ToLower(feature) {
	case "deployments":
		return deployments.NewHandler(stdout, stderr, input)
	case "environments":
		return environments.NewHandler(stdout, stderr, input)
	default:
		return nil, fmt.Errorf("feature %q does not exists, available ones are: %v", feature, AvailableFeatures)
	}
}
