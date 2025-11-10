package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cycloidio/gitlab-resource/features"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:          "gitlab-resource [in/out/check] < input.json",
		Short:        "This is a Concource CI Resource for interacting with Gitlab.",
		Long:         "See Readme.md at https://github.com/cycloidio/gitlab-resource",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// The final bin will be executed using the shebang of files located here:
			// /opt/resource/check
			// /opt/resource/in
			// /opt/resource/out
			//
			// This logic will determine the action based
			// on the file name, and forward the args to the related command.
			// We strip the first arg since the shebang sends also the filename as
			// first arg
			//
			// This trick allow to make a very small docker image, with only one binary
			// and the scripts with shebangs.
			action := filepath.Base(os.Args[1])
			switch action {
			case "check":
				return check(cmd, args[1:])
			case "in":
				return in(cmd, args[1:])
			case "out":
				return out(cmd, args[1:])
			default:
				return fmt.Errorf("invalid command %q, only check, in and out are allowed", action)
			}
		},
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:  "check",
			RunE: check,
		},
		&cobra.Command{
			Use:  "in",
			RunE: in,
			Args: cobra.ExactArgs(2),
		},
		&cobra.Command{
			Use:  "out",
			RunE: out,
			Args: cobra.ExactArgs(2),
		},
	)

	err := rootCmd.Execute()
	if err != nil {
		// rootCmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func check(cmd *cobra.Command, args []string) error {
	var input *models.Inputs
	err := internal.ReadSourceFromStdin(cmd, &input)
	if err != nil {
		return err
	}

	handler, err := features.NewFeatureHandler(input.Source.Feature)
	if err != nil {
		return err
	}

	err = handler.Validate(input)
	if err != nil {
		return fmt.Errorf("resource source validation failed: %w", err)
	}

	versions, err := handler.Check(input)
	if err != nil {
		return err
	}

	err = internal.PrintJSON(cmd.OutOrStdout(), versions)
	if err != nil {
		return fmt.Errorf("failed to output result to stdout: %w", err)
	}

	return nil
}

func in(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing out directory as first argument")
	}
	outDir := args[0]

	var input models.Inputs
	err := internal.ReadSourceFromStdin(cmd, &input)
	if err != nil {
		return err
	}

	handler, err := features.NewFeatureHandler(input.Source.Feature)
	if err != nil {
		return err
	}

	output, err := handler.In(&input, outDir)
	if err != nil {
		return err
	}

	err = internal.PrintJSON(cmd.OutOrStdout(), output)
	if err != nil {
		return fmt.Errorf("failed to output result to stdout: %w", err)
	}

	return nil
}

func out(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing out directory as first argument")
	}
	outDir := args[0]

	var input models.Inputs
	err := internal.ReadSourceFromStdin(cmd, input)
	if err != nil {
		return err
	}

	handler, err := features.NewFeatureHandler(input.Source.Feature)
	if err != nil {
		return err
	}

	output, err := handler.Out(&input, outDir)
	if err != nil {
		return err
	}

	err = internal.PrintJSON(cmd.OutOrStdout(), output)
	if err != nil {
		return fmt.Errorf("failed to output result to stdout: %w", err)
	}

	return nil
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
