package env

import "os"

// Repository ...
type Repository interface {
	List() []string
	Unset(key string) error
	Get(key string) string
	Set(key, value string) error
}

// NewRepository ...
func NewRepository() Repository {
	return defaultRepository{}
}

type defaultRepository struct{}

// Get ...
func (d defaultRepository) Get(key string) string {
	return os.Getenv(key)
}

// Set ...
func (d defaultRepository) Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset ...
func (d defaultRepository) Unset(key string) error {
	return os.Unsetenv(key)
}

// List ...
func (d defaultRepository) List() []string {
	return os.Environ()
}
