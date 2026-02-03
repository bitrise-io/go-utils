package progress

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Ticker provides periodic output for long-running operations in CI environments.
// It prints dots at regular intervals to show that work is progressing.
type Ticker struct {
	message  string
	interval time.Duration
	writer   io.Writer
	sleeper  Sleeper

	mu       sync.Mutex
	active   bool
	stopChan chan bool
}

// NewTicker creates a new Ticker with the given parameters.
func NewTicker(message string, interval time.Duration, writer io.Writer) *Ticker {
	return NewTickerWithSleeper(message, interval, writer, DefaultSleeper{})
}

// NewTickerWithSleeper creates a new Ticker with a custom Sleeper for testing.
func NewTickerWithSleeper(message string, interval time.Duration, writer io.Writer, sleeper Sleeper) *Ticker {
	return &Ticker{
		message:  message,
		interval: interval,
		writer:   writer,
		sleeper:  sleeper,

		active:   false,
		stopChan: make(chan bool),
	}
}

// NewDefaultTicker creates a Ticker with a default 5-second interval.
func NewDefaultTicker(message string, writer io.Writer) *Ticker {
	return NewTicker(message, 5*time.Second, writer)
}

// Start begins the ticker. Prints the initial message and starts periodic dot output.
func (t *Ticker) Start() {
	t.mu.Lock()
	if t.active {
		t.mu.Unlock()
		return
	}
	t.active = true
	t.mu.Unlock()

	// Print initial message
	if _, err := fmt.Fprint(t.writer, t.message); err != nil {
		fmt.Printf("failed to print message: %s, error: %s\n", t.message, err)
	}

	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			default:
				t.sleeper.Sleep(t.interval)

				t.mu.Lock()
				if !t.active {
					t.mu.Unlock()
					return
				}
				if _, err := fmt.Fprint(t.writer, "."); err != nil {
					fmt.Printf("failed to print progress dot, error: %s\n", err)
				}
				t.mu.Unlock()
			}
		}
	}()
}

// Stop stops the ticker and prints a newline.
func (t *Ticker) Stop() {
	t.mu.Lock()
	if !t.active {
		t.mu.Unlock()
		return
	}
	t.active = false
	t.mu.Unlock()

	// Send stop signal (non-blocking - goroutine may have already exited)
	select {
	case t.stopChan <- true:
	default:
	}

	// Print newline to finish the line
	if _, err := fmt.Fprintln(t.writer); err != nil {
		fmt.Printf("failed to print newline, error: %s\n", err)
	}
}
