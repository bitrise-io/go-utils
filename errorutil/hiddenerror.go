package errorutil

import "fmt"

// HiddenOriginalError allows to include an error in the error chain but do not print it.
// this allows replacing it with a more readable error message,
// while allowing code to check for the type of the error
type HiddenOriginalError struct {
	originalErr error
}

// NewHiddenOriginalError ...
func NewHiddenOriginalError(originalErr error) *HiddenOriginalError {
	return &HiddenOriginalError{
		originalErr: originalErr,
	}
}

// Error ...
func (h HiddenOriginalError) Error() string {
	return fmt.Sprintf("%T reworded", h.originalErr)
}

// Unwrap ...
func (h HiddenOriginalError) Unwrap() error {
	return h.originalErr
}
