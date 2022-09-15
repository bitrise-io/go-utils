package logger

import (
	"strings"
	"unicode"
)

const resetCode = "\u001b[0m"

var ansiEscapeCodeToLevel = map[string]Level{
	"\u001b[31;1m": ErrorLevel,
	"\u001b[33;1m": WarnLevel,
	"\u001b[34;1m": InfoLevel,
	"\u001b[32;1m": DoneLevel,
	"\u001b[35;1m": DebugLevel,
}

func convertColoredString(message string) (Level, string) {
	logLevel := NormalLevel

	for code, level := range ansiEscapeCodeToLevel {
		if strings.HasPrefix(message, code) {
			logLevel = level
			message = strings.TrimPrefix(message, code)
			hasNewline := strings.HasSuffix(message, "\n")
			// We need to remove all the possible noise from the end as we need remove the reset ansi code from the end
			message = strings.TrimRightFunc(message, unicode.IsSpace)

			if strings.HasSuffix(message, resetCode) {
				message = strings.TrimSuffix(message, resetCode)
			}

			if hasNewline {
				message += "\n"
			}

			break
		}
	}

	return logLevel, message
}
