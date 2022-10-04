package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_converterConversion(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		expectedLevel   Level
		expectedMessage string
	}{
		{
			name:            "Normal message without a color literal",
			message:         "This is a normal message without a color literal\n",
			expectedLevel:   NormalLevel,
			expectedMessage: "This is a normal message without a color literal\n",
		},
		{
			name:            "Error message",
			message:         "\u001B[31;1mThis is an error\u001B[0m",
			expectedLevel:   ErrorLevel,
			expectedMessage: "This is an error",
		},
		{
			name:            "Warn message",
			message:         "\u001B[33;1mThis is a warning\u001B[0m",
			expectedLevel:   WarnLevel,
			expectedMessage: "This is a warning",
		},
		{
			name:            "Info message",
			message:         "\u001B[34;1mThis is an Info\u001B[0m",
			expectedLevel:   InfoLevel,
			expectedMessage: "This is an Info",
		},
		{
			name:            "Done message",
			message:         "\u001B[32;1mThis is a done message\u001B[0m",
			expectedLevel:   DoneLevel,
			expectedMessage: "This is a done message",
		},
		{
			name:            "Debug message",
			message:         "\u001B[35;1mThis is a debug message\u001B[0m",
			expectedLevel:   DebugLevel,
			expectedMessage: "This is a debug message",
		},
		{
			name:            "Error message with whitespaces at the end",
			message:         "\u001B[31;1mLast error\u001B[0m   \n",
			expectedLevel:   ErrorLevel,
			expectedMessage: "Last error   \n",
		},
		{
			name:            "Error message with whitespaces at the beginning",
			message:         "  \u001B[31;1mLast error\u001B[0m   \n",
			expectedLevel:   NormalLevel,
			expectedMessage: "  \u001B[31;1mLast error\u001B[0m   \n",
		},
		{
			name:            "Error message without a closing color literal",
			message:         "\u001B[31;1mAnother error\n",
			expectedLevel:   NormalLevel,
			expectedMessage: "\u001B[31;1mAnother error\n",
		},
		{
			name:            "Info message with multiple embedded colors",
			message:         "\u001B[34;1mThis is \u001B[33;1mmulti color \u001B[31;1mInfo message\u001B[0m",
			expectedLevel:   NormalLevel,
			expectedMessage: "\u001B[34;1mThis is \u001B[33;1mmulti color \u001B[31;1mInfo message\u001B[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, message := convertColoredString(tt.message)

			assert.Equal(t, tt.expectedLevel, level)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}
