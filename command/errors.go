package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ExitStatusError ...
type ExitStatusError struct {
	readableReason  error
	originalExitErr error
}

// NewExitStatusError ...
func NewExitStatusError(printableCmdArgs string, exitErr *exec.ExitError, errorLines []string) error {
	reasonMsg := fmt.Sprintf("command failed with exit status %d (%s)", exitErr.ExitCode(), printableCmdArgs)

	errorOutput := strings.Join(errorLines, "\n")
	if len(errorOutput) == 0 {
		if len(exitErr.Stderr) != 0 {
			errorOutput = string(exitErr.Stderr)
		} else {
			errorOutput = "check the command's output for details"
		}
	}

	return &ExitStatusError{
		readableReason:  fmt.Errorf("%s: %w", reasonMsg, errors.New(errorOutput)),
		originalExitErr: exitErr,
	}
}

// Error returns the reason as a single line. The places that show a command's failure in one line cut
// the message at its first newline, which would drop every reason line but the first. The layered form
// is in Reason. Does not include the original error message (`exit status 1`).
func (e *ExitStatusError) Error() string {
	var parts []string
	for _, line := range strings.Split(e.readableReason.Error(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts = append(parts, line)
	}

	return strings.Join(parts, " ")
}

// Unwrap is needed for errors.Is and errors.As to work correctly.
func (e *ExitStatusError) Unwrap() error {
	return e.originalExitErr
}

// Reason returns the user-friendly error, to be used by errorutil.ErrorFormatter.
func (e *ExitStatusError) Reason() error {
	return e.readableReason
}
