package errorutil

import (
	"errors"
)

// DebugError type errors will not be included in FormattedError output
// but due to Unwrap() any error that is wrapped by an InternalDebugError can be used with errors.As()
type DebugError interface {
	Error() string
	Unwrap() error
	OriginalError() error
}

// JoinInternalDebugError ...
func JoinInternalDebugError(err error, debugErr error) error {
	return errors.Join(err, NewInternalDebugError(debugErr))
}

// InternalDebugError allows to include an error in the error chain but do not print it.
// this allows replacing it with a more readable error message,
// while allowing code to check for the type of the error
type InternalDebugError struct {
	originalErr error
}

// NewInternalDebugError ...
func NewInternalDebugError(originalErr error) error {
	return &InternalDebugError{
		originalErr: originalErr,
	}
}

// Error ...
func (e *InternalDebugError) Error() string {
	return "" // do not print this error as it is internal debug info
}

// Unwrap ...
func (e *InternalDebugError) Unwrap() error {
	return e.originalErr
}

// OriginalError ...
func (e *InternalDebugError) OriginalError() error {
	return e.originalErr
}
