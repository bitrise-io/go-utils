package os

import "os"

// EnvironmentRepository ...
type EnvironmentRepository interface {
	SetVariable(key string, value string) error
	UnsetVariable(key string) error
	ListVariable() []string
}

type defaultEnvironmentRepository struct{}

// NewEnvironmentRepository ...
func NewEnvironmentRepository() EnvironmentRepository {
	return defaultEnvironmentRepository{}
}

// SetVariable ...
func (m defaultEnvironmentRepository) SetVariable(key string, value string) error {
	return os.Setenv(key, value)
}

// UnsetVariable ...
func (m defaultEnvironmentRepository) UnsetVariable(key string) error {
	return os.Unsetenv(key)
}

// ListVariable ...
func (m defaultEnvironmentRepository) ListVariable() []string {
	return os.Environ()
}
