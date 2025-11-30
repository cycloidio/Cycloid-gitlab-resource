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

// Feature interface defines the three actions that each reature must implement
// check, in or out.
// The inputs are passed in the NewFeatureHandler call that each feature must
// deserialize with its own types.
type Feature interface {
	Check() error
	In(outDir string) error
	Out(outDir string) error
}

// NewFeatureHandler will output a the correct handler depending on the feature
// It accepts the input verbatim as []byte so that each handler can serialize with
// the correct type.
// We also git a reference to stdout and stderr for handler to use to output informations.
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
