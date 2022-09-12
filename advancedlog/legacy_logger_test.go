package logger

import (
	"testing"

	"github.com/bitrise-io/go-utils/v2/log/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_GivenLEgacyLogger_WhenLogMessageInvoked_ThenLogsItCorrectly(t *testing.T) {
	tests := []struct {
		name                string
		enableDebugLogs     bool
		hasOutput           bool
		parameters          testLogParameters
		expectedLogFunction string
		expectedMessage     string
	}{
		{
			name:            "Error log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    ErrorLevel,
				message:  "Error",
			},
			expectedLogFunction: "Errorf",
			expectedMessage:     "Error",
		},
		{
			name:            "Warning log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    WarnLevel,
				message:  "Warning",
			},
			expectedLogFunction: "Warnf",
			expectedMessage:     "Warning",
		},
		{
			name:            "Info log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: CLI,
				level:    InfoLevel,
				message:  "Info",
			},
			expectedLogFunction: "Infof",
			expectedMessage:     "Info",
		},
		{
			name:            "Done log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: CLI,
				level:    DoneLevel,
				message:  "Done",
			},
			expectedLogFunction: "Donef",
			expectedMessage:     "Done",
		},
		{
			name:            "Normal log",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    NormalLevel,
				message:  "Normal",
			},
			expectedLogFunction: "Printf",
			expectedMessage:     "Normal",
		},
		{
			name:            "Debug log",
			enableDebugLogs: true,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: Step,
				level:    DebugLevel,
				message:  "Debug",
			},
			expectedLogFunction: "Debugf",
			expectedMessage:     "Debug",
		},
		{
			name:            "Debug log is not logged when disabled",
			enableDebugLogs: false,
			hasOutput:       false,
			parameters: testLogParameters{
				producer: Step,
				level:    DebugLevel,
				message:  "Debug",
			},
			expectedLogFunction: "Debugf",
			expectedMessage:     "Debug",
		},
		{
			name:            "Empty message is logged",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: CLI,
				level:    InfoLevel,
				message:  "\n",
			},
			expectedLogFunction: "Infof",
			expectedMessage:     "\n",
		},
		{
			name:            "Closing newline is removed from the message",
			enableDebugLogs: false,
			hasOutput:       true,
			parameters: testLogParameters{
				producer: CLI,
				level:    InfoLevel,
				message:  "This is the message\n",
			},
			expectedLogFunction: "Infof",
			expectedMessage:     "This is the message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &mocks.Logger{}
			mockLogger.On(tt.expectedLogFunction, mock.Anything).Return()
			mockLogger.On("EnableDebugLog", mock.Anything).Return()

			logger := newLegacyLogger(mockLogger)
			logger.EnableDebugLog(tt.enableDebugLogs)
			logger.LogMessage(tt.parameters.producer, tt.parameters.level, tt.parameters.message)

			mockLogger.AssertCalled(t, "EnableDebugLog", tt.enableDebugLogs)

			if tt.hasOutput {
				mockLogger.AssertCalled(t, tt.expectedLogFunction, tt.expectedMessage)
			} else {
				mockLogger.AssertNotCalled(t, tt.expectedLogFunction, tt.expectedMessage)
			}
		})
	}
}
