package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// DetectStdinInput will return true if the current
// command effectively has some kind of stdin input
// return false otherwise
func DetectStdinInput() bool {
	stats, _ := os.Stdin.Stat()
	return (stats.Mode() & os.ModeCharDevice) == 0
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

func GetUser(userID int, client *gitlab.Client) (*gitlab.User, error) {
	user, _, err := client.Users.GetUser(userID, gitlab.GetUsersOptions{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user email with id %d: %w", userID, err)
	}

	return user, nil
}
