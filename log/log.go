package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger ...
type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Printf(format string, v ...interface{})
	Donef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Println()
}

type commandLineLogger struct {
}

// NewCommandLineLogger ...
func NewCommandLineLogger() Logger {
	return commandLineLogger{}
}

// Infof ...
func (l commandLineLogger) Infof(format string, v ...interface{}) {
	Infof(format, v...)
}

// Warnf ...
func (l commandLineLogger) Warnf(format string, v ...interface{}) {
	Warnf(format, v...)
}

// Printf ...
func (l commandLineLogger) Printf(format string, v ...interface{}) {
	Printf(format, v...)
}

// Donef ...
func (l commandLineLogger) Donef(format string, v ...interface{}) {
	Donef(format, v...)
}

// Debugf ...
func (l commandLineLogger) Debugf(format string, v ...interface{}) {
	Debugf(format, v...)
}

// Errorf ...
func (l commandLineLogger) Errorf(format string, v ...interface{}) {
	Errorf(format, v...)
}

// Println ...
func (l commandLineLogger) Println() {
	fmt.Println()
}

var outWriter io.Writer = os.Stdout

// SetOutWriter ...
func SetOutWriter(writer io.Writer) {
	outWriter = writer
}

var enableDebugLog = false

// SetEnableDebugLog ...
func SetEnableDebugLog(enable bool) {
	enableDebugLog = enable
}

var timestampLayout = "15:04:05"

// SetTimestampLayout ...
func SetTimestampLayout(layout string) {
	timestampLayout = layout
}

func timestampField() string {
	currentTime := time.Now()
	return fmt.Sprintf("[%s]", currentTime.Format(timestampLayout))
}
