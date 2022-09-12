package logger

import "strings"

const resetCode = "\u001b[0m"

var mapping = map[string]Level{
	"\u001b[31;1m": ErrorLevel,
	"\u001b[33;1m": WarnLevel,
	"\u001b[34;1m": InfoLevel,
	"\u001b[32;1m": DoneLevel,
	resetCode:      NormalLevel,
	"\u001b[35;1m": DebugLevel,
}

func convertColoredString(message string) (Level, string) {
	logLevel := NormalLevel

	for code, level := range mapping {
		if strings.HasPrefix(message, code) {
			logLevel = level
			message = strings.TrimPrefix(message, code)

			if strings.HasSuffix(message, resetCode) {
				message = strings.TrimSuffix(message, resetCode)
			} else if strings.HasSuffix(message, resetCode+"\n") {
				message = strings.TrimSuffix(message, resetCode+"\n")
				message += "\n"
			}
		}
	}

	return logLevel, message
}
