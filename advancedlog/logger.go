package logger

import "os"

const (
	// OutputFormatKey ...
	OutputFormatKey = "output-format"
	// JSONFormat ...
	JSONFormat = "json"
)

// Logger ...
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Done(args ...interface{})
	Donef(format string, args ...interface{})
	Doneln(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	TPrintf(format string, args ...interface{})
	EnableDebugLog(enable bool)
	IsDebugLogEnabled() bool
}

// SimplifiedLogger ...
type SimplifiedLogger interface {
	EnableDebugLog(enabled bool)
	IsDebugLogEnabled() bool
	LogMessage(producer Producer, level Level, message string)
}

// DefaultLogger ...
var DefaultLogger = newMainLogger()

// SetupLogger ...
func SetupLogger(outputFormat string) {
	if outputFormat == JSONFormat {
		DefaultLogger.setInternalLogger(newJSONLogger(os.Stdout, defaultTimeProvider))
	}
}

// EnableDebugLog ...
func EnableDebugLog(enable bool) {
	DefaultLogger.EnableDebugLog(enable)
}
