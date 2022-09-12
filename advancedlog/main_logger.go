package logger

import (
	"fmt"
	
	logutils "github.com/bitrise-io/go-utils/v2/log"
)

type mainLogger struct {
	internalLogger  SimplifiedLogger
	debugLogEnabled bool
}

func (m *mainLogger) Debug(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DebugLevel, fmt.Sprint(args...))
}

func (m *mainLogger) Debugf(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DebugLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Debugln(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DebugLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Info(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, InfoLevel, fmt.Sprint(args...))
}

func (m *mainLogger) Infof(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, InfoLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Infoln(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, InfoLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Done(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DoneLevel, fmt.Sprint(args...))
}

func (m *mainLogger) Donef(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DoneLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Doneln(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, DoneLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Warn(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, WarnLevel, fmt.Sprint(args...))
}

func (m *mainLogger) Warnf(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, WarnLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Warnln(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, WarnLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Error(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, ErrorLevel, fmt.Sprint(args...))
}

func (m *mainLogger) Errorf(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, ErrorLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Errorln(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, ErrorLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Fatal(args ...interface{}) {
	m.Error(args...)
}

func (m *mainLogger) Fatalf(format string, args ...interface{}) {
	m.Errorf(format, args...)
}

func (m *mainLogger) Fatalln(args ...interface{}) {
	m.Errorln(args...)
}

func (m *mainLogger) Print(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, NormalLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) Printf(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, NormalLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) Println(args ...interface{}) {
	m.internalLogger.LogMessage(CLI, NormalLevel, fmt.Sprintln(args...))
}

func (m *mainLogger) TPrintf(format string, args ...interface{}) {
	m.internalLogger.LogMessage(CLI, NormalLevel, fmt.Sprintf(format, args...))
}

func (m *mainLogger) EnableDebugLog(enable bool) {
	m.internalLogger.EnableDebugLog(enable)
	m.debugLogEnabled = enable
}

func (m *mainLogger) IsDebugLogEnabled() bool {
	return m.debugLogEnabled
}

func (m *mainLogger) setInternalLogger(logger SimplifiedLogger) {
	m.internalLogger = logger
	m.internalLogger.EnableDebugLog(m.debugLogEnabled)
}

func (m *mainLogger) LogMessage(producer Producer, level Level, message string) {
	m.internalLogger.LogMessage(producer, level, message)
}

func newMainLogger() mainLogger {
	return mainLogger{
		internalLogger:  newLegacyLogger(logutils.NewLogger()),
		debugLogEnabled: false,
	}
}
