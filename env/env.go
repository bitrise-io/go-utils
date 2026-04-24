package env

import (
	"fmt"
	"os"
	"os/exec"
)

// CommandLocator ...
type CommandLocator interface {
	LookPath(file string) (string, error)
}

type commandLocator struct{}

// NewCommandLocator ...
func NewCommandLocator() CommandLocator {
	return commandLocator{}
}

// LookPath ...
func (l commandLocator) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Repository abstracts read/write access to process environment variables.
// Implementations should be safe to replace in tests (e.g. with an
// in-memory fake) without touching the real process environment.
type Repository interface {
	// List returns the current environment as "KEY=VALUE" entries, like
	// os.Environ().
	List() []string
	// Unset removes key from the environment.
	Unset(key string) error
	// Get returns the value of key, or "" when unset.
	Get(key string) string
	// Set assigns value to key.
	Set(key, value string) error

	// GetOrDefault returns the value of key or def when the key is unset
	// or empty.
	GetOrDefault(key, def string) string
	// Required returns the value of key or an error when it is unset or
	// empty.
	Required(key string) (string, error)
	// FlagOrEnv returns *flag when it points to a non-empty string, else
	// the value of key. A nil pointer is treated as unset.
	FlagOrEnv(flag *string, key string) string

	// Revokable sets key to value and returns a function that restores
	// the previous value when invoked.
	Revokable(key, value string) (revoke func() error, err error)
	// RevokableMany applies every entry in envs and returns a single
	// function that restores all previous values. If any Set fails,
	// the returned revoke still restores entries already applied.
	RevokableMany(envs map[string]string) (revoke func() error, err error)
	// Scoped sets key to value, invokes fn, and then restores the
	// previous value. The restore runs even if fn panics.
	Scoped(key, value string, fn func()) error
	// ScopedMany applies every entry in envs, invokes fn, and restores
	// all previous values. Restore runs even if fn panics.
	ScopedMany(envs map[string]string, fn func()) error
}

// NewRepository ...
func NewRepository() Repository {
	return repository{}
}

type repository struct{}

// Get ...
func (d repository) Get(key string) string {
	return os.Getenv(key)
}

// Set ...
func (d repository) Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset ...
func (d repository) Unset(key string) error {
	return os.Unsetenv(key)
}

// List ...
func (d repository) List() []string {
	return os.Environ()
}

// GetOrDefault ...
func (d repository) GetOrDefault(key, def string) string {
	if v := d.Get(key); v != "" {
		return v
	}
	return def
}

// Required ...
func (d repository) Required(key string) (string, error) {
	if v := d.Get(key); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("required environment variable (%s) not provided", key)
}

// FlagOrEnv ...
func (d repository) FlagOrEnv(flag *string, key string) string {
	if flag != nil && *flag != "" {
		return *flag
	}
	return d.Get(key)
}

// Revokable ...
func (d repository) Revokable(key, value string) (func() error, error) {
	orig := d.Get(key)
	revoke := func() error { return d.Set(key, orig) }
	return revoke, d.Set(key, value)
}

// RevokableMany sets every key in envs and returns a revoke function that restores
// the previous values. If any Set fails, every key already written is restored
// before returning; the returned error wraps both the Set failure and any
// restore failure, and the returned revoke is a no-op.
func (d repository) RevokableMany(envs map[string]string) (func() error, error) {
	originals := make(map[string]string, len(envs))
	revoke := func() error {
		for k, v := range originals {
			if err := d.Set(k, v); err != nil {
				return err
			}
		}
		return nil
	}

	for k, v := range envs {
		originals[k] = d.Get(k)
		if err := d.Set(k, v); err != nil {
			if rerr := revoke(); rerr != nil {
				return func() error { return nil }, fmt.Errorf("set %q: %w (restore failed: %v)", k, err, rerr)
			}
			return func() error { return nil }, fmt.Errorf("set %q: %w", k, err)
		}
	}
	return revoke, nil
}

// Scoped ...
func (d repository) Scoped(key, value string, fn func()) (err error) {
	revoke, setErr := d.Revokable(key, value)
	if setErr != nil {
		return setErr
	}
	defer func() {
		if rerr := revoke(); err == nil {
			err = rerr
		}
	}()
	fn()
	return nil
}

// ScopedMany ...
func (d repository) ScopedMany(envs map[string]string, fn func()) (err error) {
	revoke, setErr := d.RevokableMany(envs)
	if setErr != nil {
		return setErr
	}
	defer func() {
		if rerr := revoke(); err == nil {
			err = rerr
		}
	}()
	fn()
	return nil
}
