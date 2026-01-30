package progress

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// MockSleeper is a mock implementation of Sleeper for testing.
type MockSleeper struct {
	sleepCalls []time.Duration
}

// Sleep records the duration and doesn't actually sleep.
func (m *MockSleeper) Sleep(d time.Duration) {
	m.sleepCalls = append(m.sleepCalls, d)
}

func TestDefaultSleeper_Sleep(t *testing.T) {
	sleeper := DefaultSleeper{}
	start := time.Now()
	sleeper.Sleep(10 * time.Millisecond)
	elapsed := time.Since(start)

	require.GreaterOrEqual(t, elapsed, 10*time.Millisecond, "should sleep for at least the requested duration")
}
