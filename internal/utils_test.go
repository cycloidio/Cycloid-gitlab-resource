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
		testPayload := models.Inputs{
			Source: models.Source{
				ProjectID: "1",
				ServerURL: "https://gitlab.com",
				Deployments: &models.SourceDeployment{
					Environment: Ptr("test"),
				},
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: map[string]string{"version": "1"},
		}
		cmd.SetIn(strings.NewReader(string(TMustJSON(t, testPayload))))

		var input models.Inputs
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})

	t.Run("InInputOk", func(t *testing.T) {
		testPayload := models.Inputs{
			Source: models.Source{
				ProjectID: "1",
				ServerURL: "https://gitlab.com",
				Deployments: &models.SourceDeployment{
					Environment: Ptr("test"),
				},
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: map[string]string{"version": "1"},
			Params:  &models.Params{},
		}

		cmd.SetIn(
			strings.NewReader(string(TMustJSON(t, testPayload))),
		)

		var input models.Inputs
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})

	t.Run("OutInputOk", func(t *testing.T) {
		testPayload := models.Inputs{
			Source: models.Source{
				ProjectID: "1",
				ServerURL: "https://gitlab.com",
				Deployments: &models.SourceDeployment{
					Environment: Ptr("test"),
				},
				Auth: models.SourceAuthentication{
					Token: "osef",
				},
			},
			Version: map[string]string{"version": "1"},
			Params:  &models.Params{},
		}

		cmd.SetIn(
			strings.NewReader(string(TMustJSON(t, testPayload))),
		)

		var input models.Inputs
		err := ReadSourceFromStdin(&cmd, input)
		assert.NoError(t, err, "reading input from stdin should not err")
	})
}
