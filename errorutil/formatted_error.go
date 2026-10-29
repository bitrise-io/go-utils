package errorutil

import (
	"errors"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
)

// FormattedError ...
func FormattedError(err error) string {
	var formatted string

	i := -1
	for {
		i++

		// Use the user-friendly error message, ignore the original exec.ExitError.
		if commandExitStatusError, ok := err.(*command.ExitStatusError); ok {
			err = commandExitStatusError.Reason()
		}

		reason := err.Error()
		if err = errors.Unwrap(err); err == nil {
			formatted = appendError(formatted, reason, i, true)
			return formatted
		}

		reason = strings.TrimSuffix(reason, err.Error())
		reason = strings.TrimRight(reason, " ")
		reason = strings.TrimSuffix(reason, ":")

		formatted = appendError(formatted, reason, i, false)
	}
}

func appendError(errorMessage, reason string, i int, last bool) string {
	if i == 0 {
		errorMessage = indentedReason(reason, i)
	} else {
		errorMessage += "\n"
		errorMessage += indentedReason(reason, i)
	}

	if !last {
		errorMessage += ":"
	}

	return errorMessage
}

func indentedReason(reason string, level int) string {
	prefix := strings.Repeat("  ", level)

	var indented []string
	for _, line := range reasonLines(reason) {
		indented = append(indented, prefix+line)
	}

	return strings.Join(indented, "\n")
}

// reasonLines drops the blank lines and the indentation that every line shares, and keeps what is left
// of each line's own indentation, which is how a tool says that a line is a detail of the one above it.
func reasonLines(reason string) []string {
	var lines []string
	for _, line := range strings.Split(reason, "\n") {
		line = strings.TrimRight(line, " \n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}

	shared := sharedIndent(lines)
	for i, line := range lines {
		lines[i] = line[shared:]
	}

	return lines
}

func sharedIndent(lines []string) int {
	shared := -1
	for _, line := range lines {
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if shared == -1 || indent < shared {
			shared = indent
		}
	}
	if shared < 0 {
		return 0
	}

	return shared
}
