package main

import (
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "gitlab-resource [in/out/check] < input.json",
		Short: "This is a Concource CI Resource for interacting with Gitlab.",
		Long:  "See Readme.md at https://github.com/cycloidio/gitlab-resource",
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:  "check",
			RunE: check,
			Args: cobra.NoArgs,
		},
		&cobra.Command{
			Use:  "in",
			RunE: in,
			Args: cobra.ExactArgs(1),
		},
		&cobra.Command{
			Use:  "out",
			RunE: out,
			Args: cobra.ExactArgs(1),
		},
	)
}

/*
$BUILD_ID

	The internal identifier for the build. Right now this is numeric, but it may become a UUID in the future. Treat it as an absolute reference to the build.

$BUILD_NAME

	The build number within the build's job.

$BUILD_JOB_NAME

	The name of the build's job.

$BUILD_PIPELINE_NAME

	The name of the pipeline that the build's job lives in.

$BUILD_PIPELINE_INSTANCE_VARS

	The instance vars of the instanced pipeline that the build's job lives in, serialized as JSON. See Grouping Pipelines for a definition of instanced pipelines.

$BUILD_TEAM_NAME

	The team that the build belongs to.

$BUILD_CREATED_BY

	The username that created the build. By default it is not available. See expose_build_created_by for how to opt in. This metadata field is not made available to the get step.

$ATC_EXTERNAL_URL

	The public URL for your ATC; useful for debugging.
*/

var (
	AvailableFeatures = [...]string{
		"Deployments",
	}
)

func check(cmd *cobra.Command, args []string) error {
	var input models.CheckInput
	err := internal.ReadSourceFromStdin(cmd, input)
	if err != nil {
		return err
	}

	switch input.Source.Feature {
		"Deployments":
			
	}

	return nil
}

func in(cmd *cobra.Command, args []string) error {
	var input models.InInputs
	err := internal.ReadSourceFromStdin(cmd, input)
	if err != nil {
		return err
	}

	return nil
}

func out(cmd *cobra.Command, args []string) error {
	var input models.OutInputs
	err := internal.ReadSourceFromStdin(cmd, input)
	if err != nil {
		return err
	}

	return nil
}
