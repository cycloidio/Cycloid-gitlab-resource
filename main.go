package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cycloidio/gitlab-resource/features"
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
			var action string
			if strings.HasPrefix(args[0], "/opt/resource") {
				action = filepath.Base(args[0])
			} else {
				action = args[0]
			}

			switch action {
			case "check":
				return run(action, cmd, nil)
			case "in", "out":
				if len(args) < 2 {
					return fmt.Errorf("missing out directory argument for action %q", action)
				}

				return run(action, cmd, &args[1])
			default:
				return fmt.Errorf("invalid command %q, only check, in and out are allowed", action)
			}
		},
	}

	err := rootCmd.Execute()
	if err != nil {
		// rootCmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run(action string, cmd *cobra.Command, outDir *string) error {
	stdin := cmd.InOrStdin()
	stdout := cmd.OutOrStdout()
	stderr := cmd.OutOrStdout()

	in, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read source metadata from stdin: %w", err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", string(in))

	var input models.Input
	err = json.Unmarshal(in, &input)
	if err != nil {
		return fmt.Errorf("failed to parse source data from JSON: %v: %w", string(in), err)
	}

	handler, err := features.NewFeatureHandler(stdout, stderr, input.Source.Feature, in)
	if err != nil {
		return err
	}

	switch action {
	case "check":
		return handler.Check()
	case "in":
		if outDir == nil {
			return fmt.Errorf("outDir argument is missing")
		}
		return handler.In(*outDir)
	case "out":
		if outDir == nil {
			return fmt.Errorf("outDir argument is missing")
		}

		return handler.Out(*outDir)
	default:
		return fmt.Errorf("invalid action %q, allowed ones are: check, in, out", action)
	}
}
