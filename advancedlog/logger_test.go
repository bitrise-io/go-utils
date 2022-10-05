package logger_test

import (
	"os"
	"time"

	log "github.com/bitrise-io/go-utils/v2/advancedlog"
)

func referenceTime() time.Time {
	return time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
}

func ExampleLogger() {
	var logger log.Logger

	logger = log.NewLogger(log.ConsoleLogger, log.BitriseCLI, os.Stdout, true, referenceTime)
	logger.Errorf("This is an %s", "error")

	logger = log.NewLogger(log.JSONLogger, log.BitriseCLI, os.Stdout, true, referenceTime)
	logger.Debug("This is a debug message")

	log.InitGlobalLogger(log.ConsoleLogger, log.BitriseCLI, os.Stdout, true, referenceTime)
	log.Info("This is an info message")

	// Output: [31;1mThis is an error[0m
	// {"timestamp":"2022-01-01T01:01:01Z","type":"log","producer":"bitrise_cli","level":"debug","message":"This is a debug message"}
	// [34;1mThis is an info message[0m
}
