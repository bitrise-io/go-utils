package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_producerStringValue(t *testing.T) {
	tests := []struct {
		name     string
		producer Producer
		expected string
	}{
		{
			name:     "Cli value conversion",
			producer: CLI,
			expected: "cli",
		},
		{
			name:     "Step value conversion",
			producer: Step,
			expected: "step",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.producer.String())
		})
	}
}

func Test_levelStringValue(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{
			name:     "Error level conversion",
			level:    ErrorLevel,
			expected: "error",
		},
		{
			name:     "Warn level conversion",
			level:    WarnLevel,
			expected: "warn",
		},
		{
			name:     "Info level conversion",
			level:    InfoLevel,
			expected: "info",
		},
		{
			name:     "Done level conversion",
			level:    DoneLevel,
			expected: "done",
		},
		{
			name:     "Normal level conversion",
			level:    NormalLevel,
			expected: "normal",
		},
		{
			name:     "Debug level conversion",
			level:    DebugLevel,
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}
