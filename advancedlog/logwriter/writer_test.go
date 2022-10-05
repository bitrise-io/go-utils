package logwriter_test

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/advancedlog/logwriter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func referenceTime() time.Time {
	return time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
}

func Test_GivenWriter_WhenStdoutIsUsed_ThenCapturesTheOutput(t *testing.T) {
	tests := []struct {
		name            string
		producer        logwriter.Producer
		loggerType      logwriter.LoggerType
		message         string
		expectedMessage string
	}{
		{
			name:            "ClI console log",
			producer:        logwriter.BitriseCLI,
			loggerType:      logwriter.ConsoleLogger,
			message:         "Test message",
			expectedMessage: "Test message",
		},
		{
			name:            "Step JSON log",
			producer:        logwriter.BitriseCLI,
			loggerType:      logwriter.JSONLogger,
			message:         "Test message",
			expectedMessage: `{"timestamp":"2022-01-01T01:01:01Z","type":"log","producer":"bitrise_cli","level":"normal","message":"Test message"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := logwriter.NewLogWriter(tt.loggerType, tt.producer, &buf, true, referenceTime)

			b := []byte(tt.message)

			_, err := writer.Write(b)
			assert.NoError(t, err)
			require.Equal(t, tt.expectedMessage, buf.String())
		})
	}
}

func ExampleNewLogWriter() {
	writer := logwriter.NewLogWriter(logwriter.JSONLogger, logwriter.BitriseCLI, os.Stdout, true, referenceTime)
	cmd := exec.Command("echo", "test")
	cmd.Stdout = writer
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	// Output: {"timestamp":"2022-01-01T01:01:01Z","type":"log","producer":"bitrise_cli","level":"normal","message":"test\n"}
}
