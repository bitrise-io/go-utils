package log

import (
	"fmt"
)

func printf(severity Severity, withTime bool, format string, v ...interface{}) string {
	colorFunc := severityColorFuncMap[severity]
	message := colorFunc(format, v...)
	if withTime {
		message = fmt.Sprintf("%s %s", timestampField(), message)
	}

	if _, err := fmt.Fprintln(outWriter, message); err != nil {
		fmt.Printf("failed to print message: %s, error: %s\n", message, err)
	}

	return message
}

// Successf ...
func Successf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(successSeverity, false, format, v...),
	}
}

// Donef ...
func Donef(format string, v ...interface{}) logMessage {
	return Successf(format, v...)
}

// Infof ...
func Infof(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(infoSeverity, false, format, v...),
	}
}

// Printf ...
func Printf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(normalSeverity, false, format, v...),
	}
}

// Debugf ...
func Debugf(format string, v ...interface{}) logMessage {
	if enableDebugLog {
		return logMessage{
			LogLevel: "info",
			Message: printf(debugSeverity, false, format, v...),
		}
	}

	return logMessage{
		LogLevel: "",
		Message: "",
	}
}

// Warnf ...
func Warnf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "warn",
		Message: printf(warnSeverity, false, format, v...),
	}
}

// Errorf ...
func Errorf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "error",
		Message: printf(errorSeverity, false, format, v...),
	}
}

// TSuccessf ...
func TSuccessf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(successSeverity, true, format, v...),
	}
}

// TDonef ...
func TDonef(format string, v ...interface{}) logMessage {
	return TSuccessf(format, v...)
}

// TInfof ...
func TInfof(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(infoSeverity, true, format, v...),
	}
}

// TPrintf ...
func TPrintf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "info",
		Message: printf(normalSeverity, true, format, v...),
	}
}

// TDebugf ...
func TDebugf(format string, v ...interface{}) logMessage {
	if enableDebugLog {
		return logMessage{
			LogLevel: "info",
			Message: printf(debugSeverity, true, format, v...),
		}
	}

	return logMessage{
		LogLevel: "",
		Message: "",
	}
}

// TWarnf ...
func TWarnf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "warn",
		Message: printf(warnSeverity, true, format, v...),
	}
}

// TErrorf ...
func TErrorf(format string, v ...interface{}) logMessage {
	return logMessage{
		LogLevel: "error",
		Message: printf(errorSeverity, true, format, v...),
	}
}
