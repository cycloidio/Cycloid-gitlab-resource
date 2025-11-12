package deployments

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	gitlabclient "github.com/cycloidio/gitlab-resource/clients/gitlab"
	"github.com/cycloidio/gitlab-resource/internal"
	"github.com/cycloidio/gitlab-resource/models"
)

func (h Handler) In(outDir string) error {
	versionJSON, err := json.Marshal(h.cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to serialize version to JSON: %w", err)
	}

	err = os.WriteFile(outDir+"/metadata.json", versionJSON, 0666)
	if err != nil {
		return fmt.Errorf("failed to write version to output dir %q: %w", outDir, err)
	}

	// Get the full user payload to get the user email
	client, err := gitlabclient.NewGitlabClient(&gitlabclient.GitlabConfig{
		Token: h.cfg.Source.Token,
		Url:   h.cfg.Source.ServerURL,
	})
	if err != nil {
		return err
	}

	var metdatas = make(models.Metadatas, len(h.cfg.Version))
	i := 0
	for k, v := range h.cfg.Version {
		metdatas[0] = models.Metadata{Name: k, Value: v}
		i++
	}

	if userIDStr, ok := h.cfg.Version["user_id"]; ok {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse user id %q: %w", userIDStr, err)
		}

		user, err := internal.GetUser(int(userID), client)
		if err != nil {
			return err
		}

		metdatas = append(metdatas,
			models.Metadatas{
				{Name: "user_username", Value: user.Username},
				{Name: "user_name", Value: user.Name},
				{Name: "user_email", Value: user.Email},
				{Name: "user_public_email", Value: user.PublicEmail},
			}...,
		)
	}

	output := &models.Output{
		Version:  []map[string]string{h.cfg.Version},
		Metadata: metdatas,
	}

	return OutputJSON(h.stdout, output)
}
