// Package envutil exposes the v1 envutil package-level function names on
// top of the v2 env.Repository so existing callers can migrate by import
// path rename only. New code should use env.Repository directly — the
// methods there are mockable, while these package-level wrappers touch
// the real process environment unconditionally.
package envutil

import (
	"github.com/bitrise-io/go-utils/v2/env"
)

// defaultRepo is used by every package-level helper. It is constructed
// lazily via env.NewRepository so future changes to the default
// repository implementation are picked up without touching this package.
var defaultRepo = env.NewRepository()

// RevokableSetenv sets key to value and returns a function that restores
// the previous value.
//
// Deprecated: use env.Repository.Revokable for testable env access.
func RevokableSetenv(key, value string) (func() error, error) {
	return defaultRepo.Revokable(key, value)
}

// RevokableSetenvs applies every entry in envs and returns a single
// function that restores all previous values.
//
// Deprecated: use env.Repository.RevokableMany for testable env access.
func RevokableSetenvs(envs map[string]string) (func() error, error) {
	return defaultRepo.RevokableMany(envs)
}

// SetenvForFunction sets key to value, invokes fn, and then restores the
// previous value. The restore runs even if fn panics.
//
// Deprecated: use env.Repository.Scoped for testable env access.
func SetenvForFunction(key, value string, fn func()) error {
	return defaultRepo.Scoped(key, value, fn)
}

// SetenvsForFunction applies every entry in envs, invokes fn, and then
// restores all previous values. The restore runs even if fn panics.
//
// Deprecated: use env.Repository.ScopedMany for testable env access.
func SetenvsForFunction(envs map[string]string, fn func()) error {
	return defaultRepo.ScopedMany(envs, fn)
}

// StringFlagOrEnv returns *flagValue when it points to a non-empty
// string, else the value of envKey.
//
// Deprecated: use env.Repository.FlagOrEnv for testable env access.
func StringFlagOrEnv(flagValue *string, envKey string) string {
	return defaultRepo.FlagOrEnv(flagValue, envKey)
}

// GetenvWithDefault returns the value of envKey when it is set and
// non-empty, else defValue.
//
// Deprecated: use env.Repository.GetOrDefault for testable env access.
func GetenvWithDefault(envKey, defValue string) string {
	return defaultRepo.GetOrDefault(envKey, defValue)
}

// RequiredEnv returns the value of envKey or an error when it is unset
// or empty.
//
// Deprecated: use env.Repository.Required for testable env access.
func RequiredEnv(envKey string) (string, error) {
	return defaultRepo.Required(envKey)
}
