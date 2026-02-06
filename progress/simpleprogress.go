package progress

import (
	"fmt"
	"io"
	"time"
)

// SimpleDots provides periodic output for long-running operations in CI environments.
// It prints dots at regular intervals to show that work is progressing.
type SimpleDots struct {
	writer io.Writer
	ticker Ticker

	stopChan chan bool
}

// NewDefaultSimpleDots creates a SimpleDots with a default 5-second interval.
func NewDefaultSimpleDots(writer io.Writer) *SimpleDots {
	return NewSimpleDotsWithInterval(5*time.Second, writer)
}

// NewSimpleDotsWithInterval creates a new SimpleDots with the given interval.
func NewSimpleDotsWithInterval(interval time.Duration, writer io.Writer) *SimpleDots {
	return NewSimpleDotsWithTicker(writer, NewTicker(interval))
}

// NewSimpleDotsWithTicker creates a new SimpleDots with a custom Ticker for testing.
func NewSimpleDotsWithTicker(writer io.Writer, ticker Ticker) *SimpleDots {
	return &SimpleDots{
		writer: writer,
		ticker: ticker,
	}
}

// Run starts the progress dots and executes the given action.
func (t *SimpleDots) Run(action func() error) error {
	if t.stopChan != nil {
		return fmt.Errorf("progress can only be run once")
	}
	t.stopChan = make(chan bool)
	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			case <-t.ticker.C():
				if _, err := fmt.Fprint(t.writer, "."); err != nil {
					fmt.Printf("failed to print progress dot: %s\n", err)
					return
				}
			}
		}
	}()

	actionErr := action()

	t.ticker.Stop()
	t.stopChan <- true
	// Print newline to finish the line
	if _, err := fmt.Fprintln(t.writer); err != nil {
		fmt.Printf("failed to print newline, error: %s\n", err)
	}

	return actionErr
}
