package envutil

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "environment variable exists",
			envKey:       "TEST_ENV",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "environment variable does not exist",
			envKey:       "NON_EXISTENT_ENV",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "environment variable is empty",
			envKey:       "EMPTY_ENV",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			// Test
			result := GetEnvOrDefault(tt.envKey, tt.defaultValue)

			// Assert
			if result != tt.expected {
				t.Errorf("GetEnvOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}
