package retry

import (
	"fmt"
	"time"
)

// Action ...
type Action func(attempt uint) error

// AbortableAction ...
type AbortableAction func(attempt uint) (error, bool)

// SleeperInterface is an interface for sleeping.
type SleeperInterface interface {
	Sleep(d time.Duration)
}

// Model represents the retry model configuration.
type Model struct {
	retry    uint
	waitTime time.Duration
	sleeper  SleeperInterface
}

// New creates a Model with the specified retry count and wait time.
// Uses the default time.Sleep implementation for sleeping between retries.
func New(retry uint, waitTime time.Duration) *Model {
	return NewWithSleeper(retry, waitTime, nil)
}

// NewWithSleeper creates a Model with the specified retry count, wait time, and sleeper.
// If sleeper is nil, the default time.Sleep implementation is used.
func NewWithSleeper(retry uint, waitTime time.Duration, sleeper SleeperInterface) *Model {
	if sleeper == nil {
		sleeper = DefaultSleeper{}
	}
	return &Model{
		retry:    retry,
		waitTime: waitTime,
		sleeper:  sleeper,
	}
}

// Sleeper creates a Model with only a custom sleeper.
func Sleeper(sleeper SleeperInterface) *Model {
	return NewWithSleeper(0, 0, sleeper)
}

// Times creates a Model with the specified number of retries.
func Times(retry uint) *Model {
	return NewWithSleeper(retry, 0, nil)
}

// Wait creates a Model with the specified wait time between retries.
func Wait(waitTime time.Duration) *Model {
	return NewWithSleeper(0, waitTime, nil)
}

// Times sets the number of retries on an existing Model.
func (m *Model) Times(retry uint) *Model {
	m.retry = retry
	return m
}

// Wait sets the wait time between retries on an existing Model.
func (m *Model) Wait(waitTime time.Duration) *Model {
	m.waitTime = waitTime
	return m
}

// WithSleeper sets a custom Sleeper implementation for testing purposes.
func (m *Model) WithSleeper(sleeper SleeperInterface) *Model {
	m.sleeper = sleeper
	return m
}

// Try continues executing the supplied action while this action parameter returns an error and the configured
// number of times has not been reached. Otherwise, it stops and returns the last received error.
func (m *Model) Try(action Action) error {
	return m.TryWithAbort(func(attempt uint) (error, bool) {
		return action(attempt), false
	})
}

// TryWithAbort continues executing the supplied action while this action parameter returns an error, a false bool
// value and the configured number of times has not been reached. Returning a true value from the action aborts the
// retry loop.
//
// Good for retrying actions which can return a mix of retryable and non-retryable failures.
func (m *Model) TryWithAbort(action AbortableAction) error {
	if action == nil {
		return fmt.Errorf("no action specified")
	}

	var err error
	var shouldAbort bool

	for attempt := uint(0); (0 == attempt || nil != err) && attempt <= m.retry; attempt++ {
		if attempt > 0 && m.waitTime > 0 {
			m.sleeper.Sleep(m.waitTime)
		}

		err, shouldAbort = action(attempt)

		if shouldAbort {
			break
		}
	}

	return err
}

// DefaultSleeper is the default implementation using time.Sleep.
type DefaultSleeper struct{}

// Sleep pauses the current goroutine for at least the duration d.
func (s DefaultSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
