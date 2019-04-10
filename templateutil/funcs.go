package templateutil

import (
	"os"

	"github.com/pkg/errors"
)

// Getenv returns the specified environment variable
// or an empty string if the env var key isn't set.
func Getenv(key string) string {
	return os.Getenv(key)
}

// GetenvRequired returns the specified environment variable
// or an error if the env var isn't set or if its value is empty.
func GetenvRequired(key string) (string, error) {
	if val := os.Getenv(key); len(val) > 0 {
		return val, nil
	}
	return "", errors.Errorf("no environment variable value found for key: %s", key)
}
