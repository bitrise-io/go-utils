package corelog_test

import (
	"os"
	"time"

	"github.com/bitrise-io/go-utils/v2/advancedlog/corelog"
)

func referenceTime() time.Time {
	return time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
}

func ExampleLogger() {
	var logger corelog.Logger

	logger = corelog.NewLogger(corelog.JSONLogger, os.Stdout, referenceTime)
	logger.LogMessage(corelog.BitriseCLI, corelog.DebugLevel, "Debug message")

	logger = corelog.NewLogger(corelog.RawLogger, os.Stdout, referenceTime)
	logger.LogMessage(corelog.Step, corelog.InfoLevel, "Info message")

	// Output: {"timestamp":"2022-01-01T01:01:01Z","type":"log","producer":"bitrise_cli","level":"debug","message":"Debug message"}
	// [34;1mInfo message[0m
}
