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

	reason := nestedReason(errorLines)
	if reason == nil {
		if len(exitErr.Stderr) != 0 {
			reason = errors.New(string(exitErr.Stderr))
		} else {
			reason = errors.New("check the command's output for details")
		}
	}

	return &ExitStatusError{
		readableReason:  fmt.Errorf("%s: %w", reasonMsg, reason),
		originalExitErr: exitErr,
	}
}

// nestedReason turns the collected lines into a chain of errors, one link per indentation step. A
// matcher indents a line to say that it is a detail of the line above it, and errorutil.FormattedError
// prints each link of a chain one level deeper, so the two say the same thing.
//
// A chain holds a path, not a tree: once a line is not indented further than the one before it, it joins
// the deepest link rather than starting a sibling branch that could not be printed.
func nestedReason(lines []string) error {
	var (
		links      []string
		lastIndent int
	)
	for _, line := range lines {
		text := strings.TrimLeft(line, " ")
		if text == "" {
			continue
		}

		indent := len(line) - len(text)
		if len(links) > 0 && indent <= lastIndent {
			links[len(links)-1] += "\n" + text

			continue
		}

		links = append(links, text)
		lastIndent = indent
	}

	var reason error
	for i := len(links) - 1; i >= 0; i-- {
		if reason == nil {
			reason = errors.New(links[i])

			continue
		}

		// The sentence carries on into the link below, so the punctuation that ended it does not.
		reason = fmt.Errorf("%s: %w", strings.TrimRight(links[i], " .:"), reason)
	}

	return reason
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
