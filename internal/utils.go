package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// DetectStdinInput will return true if the current
// command effectively has some kind of stdin input
// return false otherwise
func DetectStdinInput() bool {
	stats, _ := os.Stdin.Stat()
	return (stats.Mode() & os.ModeCharDevice) == 0
}

func ReadSourceFromStdin(cmd *cobra.Command, data any) error {
	if !DetectStdinInput() {
		return fmt.Errorf("did not found any stdin input, check documentation.")
	}

	stdin, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read source metadata from stdin: %w", err)
	}

	err = json.Unmarshal(stdin, data)
	if err != nil {
		return fmt.Errorf("failed to parse source data from JSON: %v: %w", string(stdin), err)
	}

	return nil
}

func PrintJSON(w io.Writer, input any) error {
	out, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to output json: %w", err)
	}

	_, err = fmt.Fprintln(w, string(out))
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	return nil
}

// MustJSON will try to output a json string from a object
// panic if there is an error.
func MustJSON(input any) string {
	out, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("failed to serialize json with: %v", input))
	}

	return string(out)
}

// Ptr return a point to the value
func Ptr[T any](t T) *T {
	return &t
}

func GetUser(projectUser *gitlab.ProjectUser, client *gitlab.Client) (*gitlab.User, error) {
	user, _, err := client.Users.GetUser(projectUser.ID, gitlab.GetUsersOptions{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user email with id %q and username %q: %w", user.ID, user.Username, err)
	}

	return user, nil
}
