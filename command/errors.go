package command

import "os/exec"

// ExitStatusError ...
type ExitStatusError struct {
	readableReason     error
	originalCommandErr error
}

// NewExitStatusError ...
func NewExitStatusError(reasonErr error, originalExitErr *exec.ExitError) error {
	if reasonErr.Error() == "" {
		panic("reason must not be empty")
	}

	return &ExitStatusError{readableReason: reasonErr, originalCommandErr: originalExitErr}
}

// Error returns the formatted error message. Does not include the original error message (`exit status 1`).
func (c *ExitStatusError) Error() string {
	return c.readableReason.Error()
}

// Unwrap is needed for errors.Is and errors.As to work correctly.
func (c *ExitStatusError) Unwrap() error {
	return c.originalCommandErr
}

// Reason returns the user-friendly error, to be used by errorutil.ErrorFormatter.
func (c *ExitStatusError) Reason() error {
	return c.readableReason
}
