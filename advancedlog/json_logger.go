package logger

import (
	"encoding/json"
	"io"
	"time"
)

// RFC3339Micro ...
const RFC3339Micro = "2006-01-02T15:04:05.999999Z07:0"

func defaultTimeProvider() time.Time {
	return time.Now()
}

type jsonLogger struct {
	debugLogEnabled bool
	encoder         *json.Encoder
	timeProvider    func() time.Time
}

func newJSONLogger(output io.Writer, provider func() time.Time) SimplifiedLogger {
	logger := jsonLogger{
		debugLogEnabled: false,
		encoder:         json.NewEncoder(output),
		timeProvider:    provider,
	}

	return &logger
}

func (j *jsonLogger) EnableDebugLog(enabled bool) {
	j.debugLogEnabled = enabled
}

func (j *jsonLogger) IsDebugLogEnabled() bool {
	return j.debugLogEnabled
}

func (j *jsonLogger) LogMessage(producer Producer, level Level, message string) {
	if j.debugLogEnabled == false && level == DebugLevel {
		return
	}

	logMessage := logMessage{
		Timestamp:   j.timeProvider().Format(RFC3339Micro),
		MessageType: "log",
		Producer:    producer.String(),
		Level:       level.String(),
		Message:     message,
	}

	err := j.encoder.Encode(logMessage)
	if err != nil {
		// This is only to satisfy errcheck.
	}
}
