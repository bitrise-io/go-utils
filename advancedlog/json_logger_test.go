package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJSONLogMessage struct {
	Timestamp   string `json:"timestamp"`
	MessageType string `json:"type"`
	Producer    string `json:"producer"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

type testLogParameters struct {
	producer Producer
	level    Level
	message  string
}

func Test_GivenJsonLogger_WhenLogMessageInvoked_ThenGeneratesCorrectMessageFormat(t *testing.T) {
	currentTime := time.Now()
	currentTimeString := currentTime.Format(RFC3339Micro)

	tests := []struct {
		name            string
		enableDebugLogs bool
		hasOutput       bool
		parameters      testLogParameters
		expectedMessage testJSONLogMessage
	}{
		{
			name:            "CLI log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: CLI,
				level:    InfoLevel,
				message:  "This is a cli log",
			},
			expectedMessage: testJSONLogMessage{
				Timestamp:   currentTimeString,
				MessageType: "log",
				Producer:    "bitrise_cli",
				Level:       "info",
				Message:     "This is a cli log",
			},
		},
		{
			name:            "Step log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    NormalLevel,
				message:  "This is a step log",
			},
			expectedMessage: testJSONLogMessage{
				Timestamp:   currentTimeString,
				MessageType: "log",
				Producer:    "step",
				Level:       "normal",
				Message:     "This is a step log",
			},
		},
		{
			name:            "Debug log",
			enableDebugLogs: true,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    DebugLevel,
				message:  "A useful debug log",
			},
			expectedMessage: testJSONLogMessage{
				Timestamp:   currentTimeString,
				MessageType: "log",
				Producer:    "step",
				Level:       "debug",
				Message:     "A useful debug log",
			},
		},
		{
			name:            "Disabled debug log",
			enableDebugLogs: false,
			hasOutput:       false,
			parameters: testLogParameters{
				producer: CLI,
				level:    DebugLevel,
				message:  "This debug log will not show up",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			logger := newJSONLogger(&buf, func() time.Time {
				return currentTime
			})
			logger.EnableDebugLog(tt.enableDebugLogs)
			logger.LogMessage(tt.parameters.producer, tt.parameters.level, tt.parameters.message)

			if tt.hasOutput {
				b, err := json.Marshal(tt.expectedMessage)
				assert.NoError(t, err)

				expected := string(b) + "\n"
				assert.Equal(t, buf.String(), expected)
			} else {
				assert.Equal(t, buf.Len(), 0)
			}
		})
	}
}

func Test_GivenJsonLogger_WhenManualErrorMessageCreation_ThenMatchesTheLogMessageFormat(t *testing.T) {
	err := fmt.Errorf("this is an error")
	currentTime := time.Now()
	currentTimeString := currentTime.Format(RFC3339Micro)

	logger := jsonLogger{
		debugLogEnabled: false,
		encoder:         json.NewEncoder(os.Stdout),
		timeProvider: func() time.Time {
			return currentTime
		},
	}

	message := logMessage{
		Timestamp:   currentTimeString,
		MessageType: "log",
		Producer:    string(CLI),
		Level:       string(ErrorLevel),
		Message:     fmt.Sprintf("log message serialization failed: %s", err),
	}
	expected, jsonErr := json.Marshal(message)
	assert.NoError(t, jsonErr)

	received := logger.logMessageForError(err)

	assert.Equal(t, string(expected), received)
}
