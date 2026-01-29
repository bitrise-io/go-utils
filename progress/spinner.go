package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

// Sleeper defines the interface for sleeping operations, allowing for dependency injection in tests.
type Sleeper interface {
	Sleep(d time.Duration)
}

// DefaultSleeper implements Sleeper using time.Sleep.
type DefaultSleeper struct{}

// Sleep calls time.Sleep with the given duration.
func (s DefaultSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Spinner displays an animated progress indicator.
type Spinner struct {
	message string
	chars   []string
	delay   time.Duration
	writer  io.Writer
	sleeper Sleeper

	mu         sync.Mutex
	active     bool
	lastOutput string
	stopChan   chan bool
}

// NewSpinner creates a new Spinner with the given parameters.
func NewSpinner(message string, chars []string, delay time.Duration, writer io.Writer) Spinner {
	return NewSpinnerWithSleeper(message, chars, delay, writer, DefaultSleeper{})
}

// NewSpinnerWithSleeper creates a new Spinner with a custom Sleeper for testing.
func NewSpinnerWithSleeper(message string, chars []string, delay time.Duration, writer io.Writer, sleeper Sleeper) Spinner {
	return Spinner{
		message: message,
		chars:   chars,
		delay:   delay,
		writer:  writer,
		sleeper: sleeper,

		active:   false,
		stopChan: make(chan bool),
	}
}

// NewDefaultSpinner creates a Spinner with default animation characters and timing, writing to stdout.
func NewDefaultSpinner(message string) Spinner {
	return NewDefaultSpinnerWithOutput(message, os.Stdout)
}

// NewDefaultSpinnerWithOutput creates a Spinner with default animation characters and timing.
func NewDefaultSpinnerWithOutput(message string, output io.Writer) Spinner {
	chars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	delay := 100 * time.Millisecond
	return NewSpinner(message, chars, delay, output)
}

// Start begins the spinner animation in a background goroutine.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			for i := 0; i < len(s.chars); i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.mu.Lock()
					s.erase()

					out := fmt.Sprintf("%s %s", s.message, s.chars[i])
					if _, err := fmt.Fprint(s.writer, out); err != nil {
						fmt.Printf("failed to update progress, error: %s\n", err)
					}
					s.lastOutput = out
					s.mu.Unlock()

					s.sleeper.Sleep(s.delay)
				}
			}
		}
	}()
}

// Stop terminates the spinner animation and clears the output.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if s.active {
		s.active = false
		s.erase()
		s.mu.Unlock()
		s.stopChan <- true
	} else {
		s.mu.Unlock()
	}
}

func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	for _, c := range []string{"\b", " ", "\b"} {
		for i := 0; i < n; i++ {
			if _, err := fmt.Fprint(s.writer, c); err != nil {
				fmt.Printf("failed to update progress, error: %s\n", err)
			}
		}
	}
	s.lastOutput = ""
}
