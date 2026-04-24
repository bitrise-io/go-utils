// Package envutil exposes the v1 envutil package-level function names on
// top of the v2 env package so existing callers can migrate by import
// path rename only. New code should call env.Revokable, env.Scoped, etc.
// directly with an env.Repository (or any value satisfying env.Getter /
// env.GetSetter), while these package-level wrappers touch the real
// process environment unconditionally.
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
// Deprecated: use env.Revokable for testable env access.
func RevokableSetenv(key, value string) (func() error, error) {
	return env.Revokable(defaultRepo, key, value)
}

// RevokableSetenvs applies every entry in envs and returns a single
// function that restores all previous values.
//
// Deprecated: use env.RevokableMany for testable env access.
func RevokableSetenvs(envs map[string]string) (func() error, error) {
	return env.RevokableMany(defaultRepo, envs)
}

// SetenvForFunction sets key to value, invokes fn, and then restores the
// previous value. The restore runs even if fn panics.
//
// Deprecated: use env.Scoped for testable env access.
func SetenvForFunction(key, value string, fn func()) error {
	return env.Scoped(defaultRepo, key, value, fn)
}

// SetenvsForFunction applies every entry in envs, invokes fn, and then
// restores all previous values. The restore runs even if fn panics.
//
// Deprecated: use env.ScopedMany for testable env access.
func SetenvsForFunction(envs map[string]string, fn func()) error {
	return env.ScopedMany(defaultRepo, envs, fn)
}

// StringFlagOrEnv returns *flagValue when it points to a non-empty
// string, else the value of envKey.
//
// Deprecated: use env.FlagOrEnv for testable env access.
func StringFlagOrEnv(flagValue *string, envKey string) string {
	return env.FlagOrEnv(defaultRepo, flagValue, envKey)
}

// GetenvWithDefault returns the value of envKey when it is set and
// non-empty, else defValue.
//
// Deprecated: use env.GetOrDefault for testable env access.
func GetenvWithDefault(envKey, defValue string) string {
	return env.GetOrDefault(defaultRepo, envKey, defValue)
}

// RequiredEnv returns the value of envKey or an error when it is unset
// or empty.
//
// Deprecated: use env.Required for testable env access.
func RequiredEnv(envKey string) (string, error) {
	return env.Required(defaultRepo, envKey)
}
