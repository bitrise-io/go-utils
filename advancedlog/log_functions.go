package logger

import "fmt"

// Debug ...
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Debugln ...
func Debugln(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, DebugLevel, fmt.Sprintln(args...))
}

// Info ...
func Info(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, InfoLevel, fmt.Sprint(args...))
}

// Infof ...
func Infof(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, InfoLevel, fmt.Sprintf(format, args...))
}

// Infoln ...
func Infoln(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, InfoLevel, fmt.Sprintln(args...))
}

// Done ...
func Done(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, DoneLevel, fmt.Sprint(args...))
}

// Donef ...
func Donef(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, DoneLevel, fmt.Sprintf(format, args...))
}

// Doneln ...
func Doneln(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, DoneLevel, fmt.Sprintln(args...))
}

// Warn ...
func Warn(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, WarnLevel, fmt.Sprint(args...))
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, WarnLevel, fmt.Sprintf(format, args...))
}

// Warnln ...
func Warnln(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, WarnLevel, fmt.Sprintln(args...))
}

// Error ...
func Error(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, ErrorLevel, fmt.Sprint(args...))
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, ErrorLevel, fmt.Sprintf(format, args...))
}

// Errorln ...
func Errorln(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, ErrorLevel, fmt.Sprintln(args...))
}

// Fatal ...
func Fatal(args ...interface{}) {
	Error(args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
}

// Fatalln ...
func Fatalln(args ...interface{}) {
	Errorln(args...)
}

// Print ...
func Print(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, NormalLevel, fmt.Sprintln(args...))
}

// Printf ...
func Printf(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, NormalLevel, fmt.Sprintf(format, args...))
}

// Println ...
func Println(args ...interface{}) {
	DefaultLogger.LogMessage(CLI, NormalLevel, fmt.Sprintln(args...))
}

// TPrintf ...
func TPrintf(format string, args ...interface{}) {
	DefaultLogger.LogMessage(CLI, NormalLevel, fmt.Sprintf(format, args...))
}

// IsDebugLogEnabled ...
func IsDebugLogEnabled() bool {
	return DefaultLogger.IsDebugLogEnabled()
}
