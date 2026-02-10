package progress

import (
	"errors"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
)

// SimpleDots provides periodic output for long-running operations in CI environments.
// It prints dots at regular intervals to show that work is progressing.
type SimpleDots struct {
	logger log.Logger
	ticker Ticker

	stopChan chan bool
}

// NewDefaultSimpleDots creates a SimpleDots with a default 5-second interval.
func NewDefaultSimpleDots(logger log.Logger) *SimpleDots {
	return NewSimpleDotsWithInterval(5*time.Second, logger)
}

// NewSimpleDotsWithInterval creates a new SimpleDots with the given interval.
func NewSimpleDotsWithInterval(interval time.Duration, logger log.Logger) *SimpleDots {
	return NewSimpleDotsWithTicker(logger, NewTicker(interval))
}

// NewSimpleDotsWithTicker creates a new SimpleDots with a custom Ticker for testing.
func NewSimpleDotsWithTicker(logger log.Logger, ticker Ticker) *SimpleDots {
	return &SimpleDots{
		logger: logger,
		ticker: ticker,
	}
}

// Run starts the progress dots and executes the given action.
func (t *SimpleDots) Run(action func() error) error {
	if t.stopChan != nil {
		return errors.New("progress can only be run once")
	}

	t.stopChan = make(chan bool)
	defer func() {
		t.ticker.Stop()
		close(t.stopChan)  // Signal the ticker goroutine to stop
		t.logger.Println() // Print a newline after the dots
	}()

	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			case <-t.ticker.C():
				t.logger.PrintWithoutNewline(".")
			}
		}
	}()

	return action()
}
