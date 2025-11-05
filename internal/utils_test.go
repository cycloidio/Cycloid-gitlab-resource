package internal

import (
	"strings"
	"testing"

	"github.com/cycloidio/gitlab-resource/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestReadSourceFromStdin(t *testing.T) {
	cmd := cobra.Command{}

	t.Run("CheckInputOk", func(t *testing.T) {
		testPayload := models.CheckInput{
			Source: models.Source{
				Project:     "test",
				ProjectID:   "1",
				Environment: "test",
				ServerURL:   "https://gitlab.com",
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: models.Version{
				map[string]string{"version": "1"},
			},
		}
		cmd.SetIn(
			strings.NewReader(string(MustJSON(t, testPayload))),
		)

		var input models.CheckInput
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})

	t.Run("InInputOk", func(t *testing.T) {
		testPayload := models.InInputs{
			Source: models.Source{
				Project:     "test",
				ProjectID:   "1",
				Environment: "test",
				ServerURL:   "https://gitlab.com",
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: models.Version{
				map[string]string{"version": "1"},
			},
			Params: models.InParams{},
		}

		cmd.SetIn(
			strings.NewReader(string(MustJSON(t, testPayload))),
		)

		var input models.CheckInput
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})

	t.Run("OutInputOk", func(t *testing.T) {
		testPayload := models.OutInputs{
			Source: models.Source{
				Project:     "test",
				ProjectID:   "1",
				Environment: "test",
				ServerURL:   "https://gitlab.com",
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: models.Version{
				map[string]string{"version": "1"},
			},
			Params: models.OutParams{},
		}

		cmd.SetIn(
			strings.NewReader(string(MustJSON(t, testPayload))),
		)

		var input models.CheckInput
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})
}
