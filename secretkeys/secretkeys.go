package secretkeys

import (
	"strings"

	"github.com/bitrise-io/go-utils/v2/env"
)

const (
	// EnvKey is the shared env var key
	EnvKey    = "BITRISE_SECRET_ENV_KEY_LIST"
	separator = ","
)

// Manager ...
type Manager interface {
	Load(envRepository env.Repository) []string
	Format(keys []string) string
}

type manager struct {
}

// NewManager creates a new instance
func NewManager() Manager {
	return manager{}
}

// Load returns the list of secret env keys
func (manager) Load(envRepository env.Repository) []string {
	value := envRepository.Get(EnvKey)
	keys := strings.Split(value, separator)
	return keys
}

// Format returns a formatted keys string
func (manager) Format(keys []string) string {
	return strings.Join(keys, separator)
}
