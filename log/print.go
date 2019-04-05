package log

import (
	"fmt"
)

func printf(severity Severity, withTime bool, format string, v ...interface{}) {
	message := createLogMsg(severity, withTime, format, v...)
	if _, err := fmt.Fprintln(outWriter, message); err != nil {
		fmt.Printf("failed to print message: %s, error: %s\n", message, err)
	}
}

func createLogMsg(severity Severity, withTime bool, format string, v ...interface{}) string {
	colorFunc := severityColorFuncMap[severity]
	message := colorFunc(format, v...)
	if withTime {
		message = prefixCurrentTime(message)
	}

	return message
}

func prefixCurrentTime(message string) string {
	return fmt.Sprintf("%s %s", timestampField(), message)
}

// Successf ...
func Successf(format string, v ...interface{}) Entry {
	printf(successSeverity, false, format, v...)
	return Entry{
		LogLevel: "info",
		Message: fmt.Sprintf(format, v...),
	}
}

// Donef ...
func Donef(format string, v ...interface{}) Entry {
	return Successf(format, v...)
}

// Infof ...
func Infof(format string, v ...interface{}) Entry {
	printf(infoSeverity, false, format, v...)
	return Entry{
		LogLevel: "info",
		Message: fmt.Sprintf(format, v...),
	}
}

// Printf ...
func Printf(format string, v ...interface{}) Entry {
	printf(normalSeverity, false, format, v...)
	return Entry{
		LogLevel: "info",
		Message: fmt.Sprintf(format, v...),
	}
}

// Debugf ...
func Debugf(format string, v ...interface{}) Entry {
	if enableDebugLog {
		printf(debugSeverity, false, format, v...)
	}

	return Entry{
		LogLevel: "info",
		Message: fmt.Sprintf(format, v...),
	}
}

// Warnf ...
func Warnf(format string, v ...interface{}) Entry {
	printf(warnSeverity, false, format, v...)
	return Entry{
		LogLevel: "warn",
		Message: fmt.Sprintf(format, v...),
	}
}

// Errorf ...
func Errorf(format string, v ...interface{}) Entry {
	printf(errorSeverity, false, format, v...)
	return Entry{
		LogLevel: "error",
		Message: fmt.Sprintf(format, v...),
	}
}

// TSuccessf ...
func TSuccessf(format string, v ...interface{}) Entry {
	printf(successSeverity, true, format, v...)
	return Entry{
		LogLevel: "info",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}

// TDonef ...
func TDonef(format string, v ...interface{}) Entry {
	return TSuccessf(format, v...)
}

// TInfof ...
func TInfof(format string, v ...interface{}) Entry {
	printf(infoSeverity, true, format, v...)
	return Entry{
		LogLevel: "info",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}

// TPrintf ...
func TPrintf(format string, v ...interface{}) Entry {
	printf(normalSeverity, true, format, v...)
	return Entry{
		LogLevel: "info",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}

// TDebugf ...
func TDebugf(format string, v ...interface{}) Entry {
	if enableDebugLog {
		printf(debugSeverity, true, format, v...)
	}

	return Entry{
		LogLevel: "info",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}

// TWarnf ...
func TWarnf(format string, v ...interface{}) Entry {
	printf(warnSeverity, true, format, v...)
	return Entry{
		LogLevel: "warn",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}

// TErrorf ...
func TErrorf(format string, v ...interface{}) Entry {
	printf(errorSeverity, true, format, v...)
	return Entry{
		LogLevel: "error",
		Message: prefixCurrentTime(fmt.Sprintf(format, v...)),
	}
}
