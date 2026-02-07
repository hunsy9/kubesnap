package version

import (
	"encoding/json"
	"net/http"
	"time"

	c "github.com/hunsy9/kubesnap/pkg/constant"
)

var Version = "dev"

type Release struct {
	TagName string `json:"tag_name"`
}

func CheckVersionUpdate() (string, error) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(c.GithubReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName != Version {
		return release.TagName, nil
	}

	return "", nil
}
