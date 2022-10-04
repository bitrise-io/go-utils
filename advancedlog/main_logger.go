package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/bitrise-io/go-utils/v2/advancedlog/corelog"
)

// MainLogger ...
type MainLogger struct {
	logger          corelog.Logger
	debugLogEnabled bool
}

// NewMainLogger ...
func NewMainLogger(t corelog.LoggerType, writer io.Writer, provider func() time.Time, debugLogEnabled bool) *MainLogger {
	coreLogger := corelog.NewLogger(t, writer, provider)
	return &MainLogger{
		logger:          coreLogger,
		debugLogEnabled: debugLogEnabled,
	}
}

// Error ...
func (m *MainLogger) Error(args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.ErrorLevel, fmt.Sprint(args...))
}

// Errorf ...
func (m *MainLogger) Errorf(format string, args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.ErrorLevel, fmt.Sprintf(format, args...))
}

// Warn ...
func (m *MainLogger) Warn(args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.WarnLevel, fmt.Sprint(args...))
}

// Warnf ...
func (m *MainLogger) Warnf(format string, args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.WarnLevel, fmt.Sprintf(format, args...))
}

// Info ...
func (m *MainLogger) Info(args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.InfoLevel, fmt.Sprint(args...))
}

// Infof ...
func (m *MainLogger) Infof(format string, args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.InfoLevel, fmt.Sprintf(format, args...))
}

// Done ...
func (m *MainLogger) Done(args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.DoneLevel, fmt.Sprint(args...))
}

// Donef ...
func (m *MainLogger) Donef(format string, args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.DoneLevel, fmt.Sprintf(format, args...))
}

// Print ...
func (m *MainLogger) Print(args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.NormalLevel, fmt.Sprintln(args...))
}

// Printf ...
func (m *MainLogger) Printf(format string, args ...interface{}) {
	m.logger.LogMessage(corelog.BitriseCLI, corelog.NormalLevel, fmt.Sprintf(format, args...))
}

// Debug ...
func (m *MainLogger) Debug(args ...interface{}) {
	if !m.debugLogEnabled {
		return
	}
	m.logger.LogMessage(corelog.BitriseCLI, corelog.DebugLevel, fmt.Sprint(args...))
}

// Debugf ...
func (m *MainLogger) Debugf(format string, args ...interface{}) {
	if !m.debugLogEnabled {
		return
	}
	m.logger.LogMessage(corelog.BitriseCLI, corelog.DebugLevel, fmt.Sprintf(format, args...))
}
