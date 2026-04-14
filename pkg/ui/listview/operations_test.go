package listview

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

// copyFile copies a file from src to dst
func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read source file %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		t.Fatalf("failed to write destination file %s: %v", dst, err)
	}
}

// getFixturePath returns absolute path to test fixture
func getFixturePath(t *testing.T, filename string) string {
	t.Helper()
	// Get the project root by finding go.mod
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	// Navigate to project root from pkg/ui/listview
	projectRoot := filepath.Join(wd, "..", "..", "..")
	fixturePath := filepath.Join(projectRoot, "test", "fixtures", filename)
	absPath, err := filepath.Abs(fixturePath)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	return absPath
}

func TestSwitchContextOperation(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFile  string
		targetCtx    string
		wantErr      bool
		errContains  string
		verifyConfig func(t *testing.T, configPath string)
	}{
		{
			name:        "switch to existing context - success",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targetCtx:   "staging-cluster",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if config.CurrentContext != "staging-cluster" {
					t.Errorf("CurrentContext = %v, want staging-cluster", config.CurrentContext)
				}
			},
		},
		{
			name:        "switch to non-existent context - error",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targetCtx:   "non-existent-cluster",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "switch to same context - success",
			fixtureFile: "kubeconfig_single_context.yaml",
			targetCtx:   "single-cluster",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if config.CurrentContext != "single-cluster" {
					t.Errorf("CurrentContext = %v, want single-cluster", config.CurrentContext)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpConfig := filepath.Join(tmpDir, "config")

			fixturePath := getFixturePath(t, tt.fixtureFile)
			copyFile(t, fixturePath, tmpConfig)

			t.Setenv("KUBECONFIG", tmpConfig)

			cmd := SwitchContextOperation(tt.targetCtx)
			msg := cmd()

			result, ok := msg.(SwitchResultMsg)
			if !ok {
				t.Fatalf("expected SwitchResultMsg, got %T", msg)
			}

			if tt.wantErr {
				if result.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !contains(result.Err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %q", result.Err, tt.errContains)
				}
			} else {
				if result.Err != nil {
					t.Errorf("unexpected error: %v", result.Err)
				}
				if tt.verifyConfig != nil {
					tt.verifyConfig(t, tmpConfig)
				}
			}
		})
	}
}

func TestSwitchNamespaceOperation(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFile  string
		targetNs     string
		wantErr      bool
		errContains  string
		verifyConfig func(t *testing.T, configPath string)
	}{
		{
			name:        "switch namespace - success",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targetNs:    "kube-system",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				ctx := config.Contexts[config.CurrentContext]
				if ctx.Namespace != "kube-system" {
					t.Errorf("Namespace = %v, want kube-system", ctx.Namespace)
				}
			},
		},
		{
			name:        "no current context set - error",
			fixtureFile: "kubeconfig_no_current.yaml",
			targetNs:    "default",
			wantErr:     true,
			errContains: "no current context",
		},
		{
			name:        "switch to default namespace - success",
			fixtureFile: "kubeconfig_single_context.yaml",
			targetNs:    "default",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				ctx := config.Contexts[config.CurrentContext]
				if ctx.Namespace != "default" {
					t.Errorf("Namespace = %v, want default", ctx.Namespace)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpConfig := filepath.Join(tmpDir, "config")

			fixturePath := getFixturePath(t, tt.fixtureFile)
			copyFile(t, fixturePath, tmpConfig)

			t.Setenv("KUBECONFIG", tmpConfig)

			cmd := SwitchNamespaceOperation(tt.targetNs)
			msg := cmd()

			result, ok := msg.(SwitchResultMsg)
			if !ok {
				t.Fatalf("expected SwitchResultMsg, got %T", msg)
			}

			if tt.wantErr {
				if result.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !contains(result.Err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %q", result.Err, tt.errContains)
				}
			} else {
				if result.Err != nil {
					t.Errorf("unexpected error: %v", result.Err)
				}
				if tt.verifyConfig != nil {
					tt.verifyConfig(t, tmpConfig)
				}
			}
		})
	}
}

func TestDeleteOperation(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFile  string
		targets      []string
		wantErr      bool
		verifyConfig func(t *testing.T, configPath string)
	}{
		{
			name:        "delete single context",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targets:     []string{"staging-cluster"},
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := config.Contexts["staging-cluster"]; exists {
					t.Error("staging-cluster should have been deleted")
				}
				if _, exists := config.Contexts["dev-cluster"]; !exists {
					t.Error("dev-cluster should still exist")
				}
				if _, exists := config.Contexts["prod-cluster"]; !exists {
					t.Error("prod-cluster should still exist")
				}
			},
		},
		{
			name:        "delete multiple contexts",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targets:     []string{"staging-cluster", "prod-cluster"},
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := config.Contexts["staging-cluster"]; exists {
					t.Error("staging-cluster should have been deleted")
				}
				if _, exists := config.Contexts["prod-cluster"]; exists {
					t.Error("prod-cluster should have been deleted")
				}
				if _, exists := config.Contexts["dev-cluster"]; !exists {
					t.Error("dev-cluster should still exist")
				}
			},
		},
		{
			name:        "delete non-existent context - no error",
			fixtureFile: "kubeconfig_single_context.yaml",
			targets:     []string{"non-existent-context"},
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := config.Contexts["single-cluster"]; !exists {
					t.Error("single-cluster should still exist")
				}
			},
		},
		{
			name:        "delete empty targets",
			fixtureFile: "kubeconfig_multi_context.yaml",
			targets:     []string{},
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				// All contexts should still exist
				if len(config.Contexts) != 3 {
					t.Errorf("expected 3 contexts, got %d", len(config.Contexts))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpConfig := filepath.Join(tmpDir, "config")

			fixturePath := getFixturePath(t, tt.fixtureFile)
			copyFile(t, fixturePath, tmpConfig)

			t.Setenv("KUBECONFIG", tmpConfig)

			cmd := DeleteOperation(tt.targets)
			msg := cmd()

			result, ok := msg.(DeletionResultMsg)
			if !ok {
				t.Fatalf("expected DeletionResultMsg, got %T", msg)
			}

			if tt.wantErr {
				if result.Err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if result.Err != nil {
					t.Errorf("unexpected error: %v", result.Err)
				}
				if tt.verifyConfig != nil {
					tt.verifyConfig(t, tmpConfig)
				}
			}
		})
	}
}

func TestRenameOperation(t *testing.T) {
	tests := []struct {
		name         string
		fixtureFile  string
		oldName      string
		newName      string
		wantErr      bool
		errContains  string
		verifyConfig func(t *testing.T, configPath string)
	}{
		{
			name:        "rename context - success",
			fixtureFile: "kubeconfig_multi_context.yaml",
			oldName:     "staging-cluster",
			newName:     "renamed-staging",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := config.Contexts["staging-cluster"]; exists {
					t.Error("old context name should not exist")
				}
				if _, exists := config.Contexts["renamed-staging"]; !exists {
					t.Error("new context name should exist")
				}
			},
		},
		{
			name:        "rename current context - updates current-context",
			fixtureFile: "kubeconfig_multi_context.yaml",
			oldName:     "dev-cluster",
			newName:     "renamed-dev",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if config.CurrentContext != "renamed-dev" {
					t.Errorf("CurrentContext = %v, want renamed-dev", config.CurrentContext)
				}
				if _, exists := config.Contexts["dev-cluster"]; exists {
					t.Error("old context name should not exist")
				}
				if _, exists := config.Contexts["renamed-dev"]; !exists {
					t.Error("new context name should exist")
				}
			},
		},
		{
			name:        "rename to existing name - error",
			fixtureFile: "kubeconfig_multi_context.yaml",
			oldName:     "staging-cluster",
			newName:     "prod-cluster",
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:        "rename non-existent context - error",
			fixtureFile: "kubeconfig_multi_context.yaml",
			oldName:     "non-existent",
			newName:     "new-name",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "rename to same name - success (no-op)",
			fixtureFile: "kubeconfig_single_context.yaml",
			oldName:     "single-cluster",
			newName:     "single-cluster",
			wantErr:     false,
			verifyConfig: func(t *testing.T, configPath string) {
				config, err := clientcmd.LoadFromFile(configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}
				if _, exists := config.Contexts["single-cluster"]; !exists {
					t.Error("context should still exist")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpConfig := filepath.Join(tmpDir, "config")

			fixturePath := getFixturePath(t, tt.fixtureFile)
			copyFile(t, fixturePath, tmpConfig)

			t.Setenv("KUBECONFIG", tmpConfig)

			// Create a mock SwitchingUIModel with the old name as choice
			model := &SwitchingUIModel{
				choice: tt.oldName,
			}

			cmd := model.RenameOperation(tt.newName)
			msg := cmd()

			result, ok := msg.(RenameResultMsg)
			if !ok {
				t.Fatalf("expected RenameResultMsg, got %T", msg)
			}

			if tt.wantErr {
				if result.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !contains(result.Err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %q", result.Err, tt.errContains)
				}
			} else {
				if result.Err != nil {
					t.Errorf("unexpected error: %v", result.Err)
				}
				if tt.verifyConfig != nil {
					tt.verifyConfig(t, tmpConfig)
				}
			}
		})
	}
}

// contains checks if substr is in s
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
