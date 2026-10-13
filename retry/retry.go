package retry

import (
	"fmt"
	"time"
)

// Action receives the number of attempts that have already failed: 0 on the first execution, 1 on
// the first retry, and so on. So `attempt > 0` means "this is a retry".
//
// Use NumberedAction with TryNumbered instead when the action wants to name the attempt it is
// running, which is the more natural thing to log.
type Action func(attempt uint) error

// AbortableAction receives the same attempt count as Action.
type AbortableAction func(attempt uint) (error, bool)

// NumberedAction receives the 1-based number of the attempt it is running: 1 on the first
// execution, 2 on the first retry, and so on.
type NumberedAction func(attempt uint) error

// Sleeper is an interface for sleeping.
type Sleeper interface {
	Sleep(d time.Duration)
}

// Model represents the retry model configuration.
type Model struct {
	retry    uint
	waitTime time.Duration
	sleeper  Sleeper
}

// New creates a Model with the specified retry count, wait time, and sleeper.
// If sleeper is nil, the default time.Sleep implementation is used.
func New(retry uint, waitTime time.Duration, sleeper Sleeper) *Model {
	if sleeper == nil {
		sleeper = DefaultSleeper{}
	}
	return &Model{
		retry:    retry,
		waitTime: waitTime,
		sleeper:  sleeper,
	}
}

// Times creates a Model with the specified number of retries.
func Times(retry uint) *Model {
	return New(retry, 0, nil)
}

// Wait creates a Model with the specified wait time between retries.
func Wait(waitTime time.Duration) *Model {
	return New(0, waitTime, nil)
}

// WithSleeper creates a Model with only a custom sleeper.
func WithSleeper(sleeper Sleeper) *Model {
	return New(0, 0, sleeper)
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
func (m *Model) WithSleeper(sleeper Sleeper) *Model {
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

// TryNumbered is Try, with the attempt numbered from 1 rather than counted from 0.
//
// It exists because the number an action wants to report is almost always the attempt it is running,
// and deriving that from Action's failure count is easy to get wrong: `attempt+1` inside an
// `attempt > 0` branch reports the attempt about to start, not the one that just failed.
//
//	err := retry.Times(2).TryNumbered(func(attempt uint) error {
//		logger.Printf("attempt %d of 3", attempt)
//		return doSomething()
//	})
//
// With TryNumbered, "this is a retry" is attempt > 1.
func (m *Model) TryNumbered(action NumberedAction) error {
	if action == nil {
		return fmt.Errorf("no action specified")
	}

	return m.TryWithAbort(func(attempt uint) (error, bool) {
		return action(attempt + 1), false
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
