package envutil

import (
	"os"
)

// GetEnvOrDefault retrieves the value of the environment variable named by the env parameter.
// If the variable is not present or is empty, it returns the defaultValue instead.
func GetEnvOrDefault(env string, defaultValue string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return defaultValue
}
