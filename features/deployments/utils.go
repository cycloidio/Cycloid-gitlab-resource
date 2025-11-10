package deployments

import (
	"encoding/json"
	"fmt"
	"io"
)

func OutputJSON(stdout io.Writer, output any) error {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize output to JSON: %w", err)
	}

	_, err = stdout.Write(data)
	if err != nil {
		return fmt.Errorf("failed to output to stdout: %w", err)
	}

	return nil
}
