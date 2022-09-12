package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testWriterParameters struct {
	producer Producer
	message  string
}

type testWriterExpectedValues struct {
	producer Producer
	level    Level
	message  string
}

func Test_GivenWriter_WhenStdoutIsUsed_ThenCapturesTheOutput(t *testing.T) {
	tests := []struct {
		name           string
		useStdout      bool
		parameters     testWriterParameters
		expectedValues testWriterExpectedValues
	}{
		{
			name:      "Cli stdout message",
			useStdout: true,
			parameters: testWriterParameters{
				producer: CLI,
				message:  "Test message",
			},
			expectedValues: testWriterExpectedValues{
				producer: CLI,
				level:    NormalLevel,
				message:  "Test message",
			},
		},
		{
			name:      "Step stderr message",
			useStdout: false,
			parameters: testWriterParameters{
				producer: Step,
				message:  "This is an error",
			},
			expectedValues: testWriterExpectedValues{
				producer: Step,
				level:    ErrorLevel,
				message:  "This is an error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedProducer Producer
			var receivedLevel Level
			var receivedMessage string

			writer := NewLogWriter(tt.parameters.producer, func(producer Producer, level Level, message string) {
				receivedProducer = producer
				receivedLevel = level
				receivedMessage = message
			})

			b := []byte(tt.parameters.message)

			if tt.useStdout {
				_, err := writer.Stdout.Write(b)
				assert.NoError(t, err)
			} else {
				_, err := writer.Stderr.Write(b)
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedValues.producer, receivedProducer)
			assert.Equal(t, tt.expectedValues.level, receivedLevel)
			assert.Equal(t, tt.expectedValues.message, receivedMessage)
		})
	}
}
