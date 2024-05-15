package errorutil

import (
	"strings"
)

func unwrap(err error) []error {
	switch x := err.(type) {
	case interface{ Unwrap() []error }:
		return x.Unwrap()
	case interface{ Unwrap() error }:
		return []error{x.Unwrap()}
	default:
		return nil
	}
}

func formattedError(err error, printDebugOnly bool, indent int) string {
	debugErr, isDebugErr := err.(DebugError)
	if isDebugErr {
		err = debugErr.OriginalError()
	}

	formatted := ""
	reason := err.Error()
	wrappedErrs := []error{}
	if wrappedErrs = unwrap(err); len(wrappedErrs) == 0 {
		if isDebugErr == printDebugOnly {
			formatted = appendError(formatted, reason, indent, true)
		}
		return formatted
	}

	for i := len(wrappedErrs) - 1; i >= 0; i-- {
		reason = strings.TrimSuffix(reason, wrappedErrs[i].Error())
		reason = strings.TrimRight(reason, "\n ")
		reason = strings.TrimSuffix(reason, ":")
	}

	if isDebugErr == printDebugOnly {
		formatted += appendError(formatted, reason, indent, false)
	}
	if !printDebugOnly && isDebugErr { // skip children of debug errors
		return formatted
	}
	if printDebugOnly && isDebugErr { // print children of debug errors
		printDebugOnly = false
	}
	for _, wrappedErr := range wrappedErrs {
		formatted += formattedError(wrappedErr, printDebugOnly, indent+1)
	}

	return formatted
}

// FormattedError ...
func FormattedError(err error) string {
	return formattedError(err, false, 0)
}

// FormattedErrorInternalDebugInfo ...
func FormattedErrorInternalDebugInfo(err error) string {
	return formattedError(err, true, 0)
}

func appendError(errorMessage, reason string, i int, last bool) string {
	if reason == "" {
		return ""
	}

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
	var lines []string
	split := strings.Split(reason, "\n")
	if len(split) == 1 && split[0] == "" {
		split = []string{"[empty error string]"}
	}

	for _, line := range split {
		line = strings.TrimLeft(line, " ")
		line = strings.TrimRight(line, "\n")
		line = strings.TrimRight(line, " ")
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	var indented string
	for i, line := range lines {
		indented += strings.Repeat("  ", level)
		indented += line
		if i != len(lines)-1 {
			indented += "\n"
		}
	}
	return indented
}
