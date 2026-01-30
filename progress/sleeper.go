package progress

import "time"

// Sleeper is an interface for sleeping, allowing for dependency injection in tests.
type Sleeper interface {
	Sleep(d time.Duration)
}

// DefaultSleeper implements Sleeper using the standard time.Sleep.
type DefaultSleeper struct{}

// Sleep calls time.Sleep with the given duration.
func (DefaultSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
