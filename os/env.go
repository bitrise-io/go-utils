package os

import "os"

// EnvironmentRepository ...
type EnvironmentRepository interface {
	SetVariable(key string, value string) error
	UnsetVariable(key string) error
	ListVariables() []string
}

type defaultEnvironmentRepository struct{}

// NewEnvironmentRepository ...
func NewEnvironmentRepository() EnvironmentRepository {
	return defaultEnvironmentRepository{}
}

// SetVariable ...
func (r defaultEnvironmentRepository) SetVariable(key string, value string) error {
	return os.Setenv(key, value)
}

// UnsetVariable ...
func (r defaultEnvironmentRepository) UnsetVariable(key string) error {
	return os.Unsetenv(key)
}

// ListVariable ...
func (r defaultEnvironmentRepository) ListVariables() []string {
	return os.Environ()
}
