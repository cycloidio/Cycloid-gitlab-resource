package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func ReadSourceFromStdin(cmd *cobra.Command, data any) error {
	stdin, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read source metadata from stdin: %w", err)
	}

	err = json.Unmarshal(stdin, &data)
	if err != nil {
		return fmt.Errorf("failed to parse source data from JSON: %v: %w", string(stdin), err)
	}

	return nil
}
