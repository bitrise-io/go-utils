package logger

import (
	logutils "github.com/bitrise-io/go-utils/v2/log"
)

type legacyLogger struct {
	debugLogEnabled bool
	logger          logutils.Logger
}

func newLegacyLogger(logger logutils.Logger) SimplifiedLogger {
	return &legacyLogger{
		debugLogEnabled: false,
		logger:          logger,
	}
}

// EnableDebugLog ...
func (l *legacyLogger) EnableDebugLog(enabled bool) {
	l.logger.EnableDebugLog(enabled)
	l.debugLogEnabled = enabled
}

// IsDebugLogEnabled ...
func (l *legacyLogger) IsDebugLogEnabled() bool {
	return l.debugLogEnabled
}

// LogMessage ...
func (l *legacyLogger) LogMessage(producer Producer, level Level, message string) {
	if !l.debugLogEnabled && level == DebugLevel {
		return
	}

	switch level {
	case ErrorLevel:
		l.logger.Errorf(message)
	case WarnLevel:
		l.logger.Warnf(message)
	case InfoLevel:
		l.logger.Infof(message)
	case DoneLevel:
		l.logger.Donef(message)
	case DebugLevel:
		l.logger.Debugf(message)
	default:
		l.logger.Printf(message)
	}
}
