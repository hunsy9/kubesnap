package version

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	c "github.com/hunsy9/kubesnap/pkg/constant"
)

func TestVersionCheckingCache_validateField(t *testing.T) {
	validTime := time.Now().Format(time.RFC3339)

	tests := []struct {
		name  string
		cache VersionCheckingCache
		want  bool
	}{
		// 성공 케이스
		{
			name: "valid cache with v prefix",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.2.3",
				LastCheckedTime: validTime,
			},
			want: true,
		},
		{
			name: "valid cache without v prefix",
			cache: VersionCheckingCache{
				LatestVersion:   "1.2.3",
				LastCheckedTime: validTime,
			},
			want: true,
		},
		{
			name: "valid dev version",
			cache: VersionCheckingCache{
				LatestVersion:   "dev",
				LastCheckedTime: validTime,
			},
			want: true,
		},

		// 실패 케이스
		{
			name: "empty last checked time",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.0.0",
				LastCheckedTime: "",
			},
			want: false,
		},
		{
			name: "invalid version format - text",
			cache: VersionCheckingCache{
				LatestVersion:   "invalid",
				LastCheckedTime: validTime,
			},
			want: false,
		},
		{
			name: "invalid version format - incomplete",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.2",
				LastCheckedTime: validTime,
			},
			want: false,
		},
		{
			name: "invalid time format",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.0.0",
				LastCheckedTime: "not-a-time",
			},
			want: false,
		},
		{
			name: "invalid time format - unix timestamp",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.0.0",
				LastCheckedTime: "1234567890",
			},
			want: false,
		},

		// 엣지 케이스
		{
			name: "empty version and time",
			cache: VersionCheckingCache{
				LatestVersion:   "",
				LastCheckedTime: "",
			},
			want: false,
		},
		{
			name: "version with extra parts",
			cache: VersionCheckingCache{
				LatestVersion:   "v1.2.3.4",
				LastCheckedTime: validTime,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cache.validateField()
			if got != tt.want {
				t.Errorf("validateField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadCache(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(tmpDir string) error
		wantEmpty      bool
		wantVersion    string
		checkDirCreate bool
	}{
		{
			name: "cache file does not exist - returns empty cache and creates directory",
			setupFunc: func(tmpDir string) error {
				// do nothing - no cache file
				return nil
			},
			wantEmpty:      true,
			checkDirCreate: true,
		},
		{
			name: "valid cache file - returns loaded cache",
			setupFunc: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return err
				}
				cache := VersionCheckingCache{
					LatestVersion:   "v1.0.0",
					LastCheckedTime: time.Now().Format(time.RFC3339),
				}
				data, _ := json.Marshal(cache)
				return os.WriteFile(filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile), data, 0644)
			},
			wantEmpty:   false,
			wantVersion: "v1.0.0",
		},
		{
			name: "invalid JSON - deletes file and returns empty cache",
			setupFunc: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile), []byte("invalid json{"), 0644)
			},
			wantEmpty: true,
		},
		{
			name: "invalid field values - deletes file and returns empty cache",
			setupFunc: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return err
				}
				cache := VersionCheckingCache{
					LatestVersion:   "invalid-version",
					LastCheckedTime: "not-a-time",
				}
				data, _ := json.Marshal(cache)
				return os.WriteFile(filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile), data, 0644)
			},
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("HOME", tmpDir)

			if err := tt.setupFunc(tmpDir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := LoadCache()

			if tt.wantEmpty {
				if got.LatestVersion != "" || got.LastCheckedTime != "" {
					t.Errorf("LoadCache() = %+v, want empty cache", got)
				}
			} else {
				if got.LatestVersion != tt.wantVersion {
					t.Errorf("LoadCache().LatestVersion = %v, want %v", got.LatestVersion, tt.wantVersion)
				}
			}

			if tt.checkDirCreate {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				if _, err := os.Stat(configDir); os.IsNotExist(err) {
					t.Errorf("LoadCache() should create config directory at %s", configDir)
				}
			}
		})
	}
}

func TestSaveCache(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		setupFunc   func(tmpDir string) error
		wantErr     bool
		verifyCache func(t *testing.T, tmpDir string)
	}{
		{
			name:    "save cache to new directory",
			version: "v1.2.3",
			setupFunc: func(tmpDir string) error {
				return nil
			},
			wantErr: false,
			verifyCache: func(t *testing.T, tmpDir string) {
				cachePath := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation, c.DefaultKubesnapVersionCacheFile)
				data, err := os.ReadFile(cachePath)
				if err != nil {
					t.Errorf("failed to read cache file: %v", err)
					return
				}
				var cache VersionCheckingCache
				if err := json.Unmarshal(data, &cache); err != nil {
					t.Errorf("failed to unmarshal cache: %v", err)
					return
				}
				if cache.LatestVersion != "v1.2.3" {
					t.Errorf("cache.LatestVersion = %v, want v1.2.3", cache.LatestVersion)
				}
				if cache.LastCheckedTime == "" {
					t.Error("cache.LastCheckedTime should not be empty")
				}
			},
		},
		{
			name:    "save cache to existing directory",
			version: "v2.0.0",
			setupFunc: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				return os.MkdirAll(configDir, 0755)
			},
			wantErr: false,
			verifyCache: func(t *testing.T, tmpDir string) {
				cachePath := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation, c.DefaultKubesnapVersionCacheFile)
				data, err := os.ReadFile(cachePath)
				if err != nil {
					t.Errorf("failed to read cache file: %v", err)
					return
				}
				var cache VersionCheckingCache
				if err := json.Unmarshal(data, &cache); err != nil {
					t.Errorf("failed to unmarshal cache: %v", err)
					return
				}
				if cache.LatestVersion != "v2.0.0" {
					t.Errorf("cache.LatestVersion = %v, want v2.0.0", cache.LatestVersion)
				}
			},
		},
		{
			name:    "overwrite existing cache",
			version: "v3.0.0",
			setupFunc: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation)
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return err
				}
				oldCache := VersionCheckingCache{
					LatestVersion:   "v1.0.0",
					LastCheckedTime: time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
				}
				data, _ := json.Marshal(oldCache)
				return os.WriteFile(filepath.Join(configDir, c.DefaultKubesnapVersionCacheFile), data, 0644)
			},
			wantErr: false,
			verifyCache: func(t *testing.T, tmpDir string) {
				cachePath := filepath.Join(tmpDir, c.DefaultKubesnapConfigLocation, c.DefaultKubesnapVersionCacheFile)
				data, err := os.ReadFile(cachePath)
				if err != nil {
					t.Errorf("failed to read cache file: %v", err)
					return
				}
				var cache VersionCheckingCache
				if err := json.Unmarshal(data, &cache); err != nil {
					t.Errorf("failed to unmarshal cache: %v", err)
					return
				}
				if cache.LatestVersion != "v3.0.0" {
					t.Errorf("cache.LatestVersion = %v, want v3.0.0", cache.LatestVersion)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("HOME", tmpDir)

			if err := tt.setupFunc(tmpDir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err := SaveCache(tt.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.verifyCache != nil {
				tt.verifyCache(t, tmpDir)
			}
		})
	}
}

func TestSaveCacheAndLoadCache_Roundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	testVersion := "v1.5.0"

	// Save cache
	if err := SaveCache(testVersion); err != nil {
		t.Fatalf("SaveCache() error = %v", err)
	}

	// Load cache
	cache := LoadCache()

	// Verify roundtrip
	if cache.LatestVersion != testVersion {
		t.Errorf("Roundtrip failed: got LatestVersion = %v, want %v", cache.LatestVersion, testVersion)
	}

	if cache.LastCheckedTime == "" {
		t.Error("Roundtrip failed: LastCheckedTime should not be empty")
	}

	// Verify LastCheckedTime is valid RFC3339
	_, err := time.Parse(time.RFC3339, cache.LastCheckedTime)
	if err != nil {
		t.Errorf("Roundtrip failed: LastCheckedTime is not valid RFC3339: %v", err)
	}

	// Verify cache is valid
	if !cache.validateField() {
		t.Error("Roundtrip failed: cache.validateField() returned false")
	}
}
