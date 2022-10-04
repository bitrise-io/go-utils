package logger

import (
	"io"
	"time"

	"github.com/bitrise-io/go-utils/v2/advancedlog/corelog"
)

// Logger ...
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Done(args ...interface{})
	Donef(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

var globalLogger Logger

func InitGlobalLogger(t corelog.LoggerType, writer io.Writer, provider func() time.Time, debugLogEnabled bool) {
	globalLogger = NewMainLogger(t, writer, provider, debugLogEnabled)
}
