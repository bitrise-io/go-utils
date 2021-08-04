package v2

import (
	"fmt"
	"github.com/bitrise-io/go-utils/log"
)

// Logger ...
type Logger interface {
	Printf(format string, v ...interface{})
	Donef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Println()
}
type defaultLogger struct {
}

// NewDefaultLogger ...
func NewDefaultLogger() Logger {
	return defaultLogger{}
}

// Printf ...
func (l defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Donef ...
func (l defaultLogger) Donef(format string, v ...interface{}) {
	log.Donef(format, v...)
}

// Debugf ...
func (l defaultLogger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Errorf ...
func (l defaultLogger) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Println ...
func (l defaultLogger) Println() {
	fmt.Println()
}
