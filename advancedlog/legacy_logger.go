package logger

import (
	logutils "github.com/bitrise-io/go-utils/v2/log"
	"strings"
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

func (l *legacyLogger) EnableDebugLog(enabled bool) {
	l.logger.EnableDebugLog(enabled)
	l.debugLogEnabled = enabled
}

func (l *legacyLogger) IsDebugLogEnabled() bool {
	return l.debugLogEnabled
}

func (l *legacyLogger) LogMessage(producer Producer, level Level, message string) {
	if l.debugLogEnabled == false && level == DebugLevel {
		return
	}

	// It is needed to trim the newline char from the end because the wrapped logger will automatically add one
	if message != "\n" {
		message = strings.TrimSuffix(message, "\n")
	}

	switch level {
	case ErrorLevel:
		l.logger.Errorf(message)
	case WarnLevel:
		l.logger.Warnf(message)
	case DoneLevel:
		l.logger.Donef(message)
	case NormalLevel:
		l.logger.Printf(message)
	case DebugLevel:
		l.logger.Debugf(message)
	default:
		l.logger.Infof(message)
	}
}
