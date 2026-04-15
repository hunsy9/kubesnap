package version

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	c "github.com/hunsy9/kubesnap/pkg/constant"
	"golang.org/x/mod/semver"
)

var Version = "dev"

type Release struct {
	TagName string `json:"tag_name"`
}

type VersionCheckingCache struct {
	LatestVersion   string `json:"latest_version"`
	LastCheckedTime string `json:"last_check_time"`
}

var versionRegex = regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

func (c *VersionCheckingCache) validateField() bool {
	if c.LastCheckedTime == "" {
		return false
	}

	if !versionRegex.MatchString(c.LatestVersion) && c.LatestVersion != "dev" {
		return false
	}

	if c.LastCheckedTime == "" {
		return false
	}

	_, err := time.Parse(time.RFC3339, c.LastCheckedTime)
	return err == nil
}

func LoadCache() VersionCheckingCache {
	home, err := os.UserHomeDir()

	if err != nil {
		return VersionCheckingCache{}
	}
	configDir := filepath.Join(home, c.DefaultKubesnapConfigLocation)
	cachePath := filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile)

	emptyCache := VersionCheckingCache{} // default cache

	data, err := os.ReadFile(cachePath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0755)
		return emptyCache
	} else if err != nil {
		return emptyCache
	}

	var cache VersionCheckingCache
	if err := json.Unmarshal(data, &cache); err != nil {
		os.Remove(cachePath)
		return emptyCache
	}

	if !cache.validateField() {
		os.Remove(cachePath)
		return emptyCache
	}

	return cache
}

func SaveCache(latestVersion string) error {
	cache := VersionCheckingCache{
		LatestVersion:   latestVersion,
		LastCheckedTime: time.Now().Format(time.RFC3339),
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, c.DefaultKubesnapConfigLocation)
	cachePath := filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// Strategy for Version Check and Cache Update:
//
// 1. If the cache is empty or invalid:
//
//   - Perform an immediate API call to fetch the latest version.
//
//     2. If the cache exists and is valid:
//     Compare the current app version with the cached 'latest version'.
//
//     Case A: Versions are Equal (App == Cached)
//
//   - Check the timestamp.
//
//   - If > 24 hours have passed since the last check: Refresh cache (API call).
//
//   - If <= 24 hours: Do nothing (skip API call).
//
//     Case B: App is Newer (App > Cached)
//
//   - This implies the user has manually upgraded the app, making the cache outdated.
//
//   - Force a cache refresh (API call) regardless of the timestamp to sync with the true latest version.
//
//     Case C: App is Older (App < Cached)
//
//   - This indicates an update is available.
//
//   - Check the timestamp (like Case A) to decide whether to refresh the cache for an even newer version.
func CheckVersionUpdate() (string, error) {

	cache := LoadCache()
	lastCheck, _ := time.Parse(time.RFC3339, cache.LastCheckedTime)

	// Skip API call if app version equals cache version and last check was within 1 day
	if cache.LatestVersion == Version && time.Since(lastCheck) < 24*time.Hour {
		return "", nil
	}

	// Call github api, and get tag(release version)
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(c.GithubReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// call github api, and get tag(release version)
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	_ = SaveCache(release.TagName)

	if semver.Compare(release.TagName, Version) > 0 || Version == "dev" {
		return release.TagName, nil
	}

	return "", nil
}
