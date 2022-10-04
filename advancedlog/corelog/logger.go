package corelog

import (
	"io"
	"time"
)

// Logger ...
type Logger interface {
	LogMessage(producer Producer, level Level, message string)
}

type LoggerType string

const (
	JSONLogger LoggerType = "json"
	RawLogger  LoggerType = "raw"
)

func NewLogger(t LoggerType, output io.Writer, provider func() time.Time) Logger {
	switch t {
	case JSONLogger:
		return newJSONLogger(output, provider)
	default:
		return newLegacyLogger(output)
	}
}
