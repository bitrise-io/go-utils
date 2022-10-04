package logger_test

import (
	"os"
	"time"

	log "github.com/bitrise-io/go-utils/v2/advancedlog"
	"github.com/bitrise-io/go-utils/v2/advancedlog/corelog"
)

func referenceTime() time.Time {
	return time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
}

func ExampleLogger() {
	var logger log.Logger

	logger = log.NewMainLogger(corelog.RawLogger, os.Stdout, referenceTime, true)
	logger.Errorf("This is an %s", "error")

	logger = log.NewMainLogger(corelog.JSONLogger, os.Stdout, referenceTime, true)
	logger.Debug("This is a debug message")

	log.InitGlobalLogger(corelog.RawLogger, os.Stdout, referenceTime, true)
	log.Info("This is an info message")

	// Output: [31;1mThis is an error[0m
	// {"timestamp":"2022-01-01T01:01:01Z","type":"log","producer":"bitrise_cli","level":"debug","message":"This is a debug message"}
	// [34;1mThis is an info message[0m
}
