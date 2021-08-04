package env

import "os"

// OsRepository ...
type OsRepository interface {
	Unset(key string) error
	Set(key string, value string) error
	List() []string
}

type osRepository struct{}

// NewOsRepository ...
func NewOsRepository() OsRepository {
	return osRepository{}
}

// List ...
func (m osRepository) List() []string {
	return os.Environ()
}

// Unset ...
func (m osRepository) Unset(key string) error {
	return os.Unsetenv(key)
}

// Set ...
func (m osRepository) Set(key string, value string) error {
	return os.Setenv(key, value)
}
