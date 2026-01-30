package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTicker(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewTicker("Testing", 1*time.Second, &buf)

	require.NotNil(t, ticker)
	require.Equal(t, "Testing", ticker.message)
	require.Equal(t, 1*time.Second, ticker.interval)
	require.Equal(t, &buf, ticker.writer)
	require.False(t, ticker.active)
}

func TestNewTickerWithSleeper(t *testing.T) {
	var buf bytes.Buffer
	sleeper := &MockSleeper{}
	ticker := NewTickerWithSleeper("Testing", 1*time.Second, &buf, sleeper)

	require.NotNil(t, ticker)
	require.Equal(t, sleeper, ticker.sleeper)
}

func TestNewDefaultTicker(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Testing", &buf)

	require.NotNil(t, ticker)
	require.Equal(t, "Testing", ticker.message)
	require.Equal(t, 5*time.Second, ticker.interval)
}

func TestTicker_StartAndStop(t *testing.T) {
	var buf bytes.Buffer
	sleeper := &MockSleeper{}
	ticker := NewTickerWithSleeper("Working", 100*time.Millisecond, &buf, sleeper)

	ticker.Start()
	require.True(t, ticker.active)

	// Initial message should be printed immediately
	require.Equal(t, "Working", buf.String())

	// Simulate some ticks
	for i := 0; i < 3; i++ {
		sleeper.Sleep(100 * time.Millisecond)
		time.Sleep(10 * time.Millisecond) // Give goroutine time to print
	}

	ticker.Stop()
	require.False(t, ticker.active)

	// Should have: "Working" + dots + newline
	output := buf.String()
	require.True(t, strings.HasPrefix(output, "Working"))
	require.True(t, strings.Contains(output, "."))
	require.True(t, strings.HasSuffix(output, "\n"))
}

func TestTicker_StartTwice(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Testing", &buf)

	ticker.Start()
	ticker.Start() // Second start should be no-op

	require.True(t, ticker.active)
	require.Equal(t, "Testing", buf.String())

	ticker.Stop()
}

func TestTicker_StopWhenNotActive(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Testing", &buf)

	// Stop without starting should be safe
	ticker.Stop()

	require.False(t, ticker.active)
	require.Equal(t, "", buf.String())
}

func TestTicker_PeriodicOutput(t *testing.T) {
	var buf bytes.Buffer
	sleeper := &MockSleeper{}
	ticker := NewTickerWithSleeper("Progress", 50*time.Millisecond, &buf, sleeper)

	ticker.Start()
	require.Equal(t, "Progress", buf.String())

	// Simulate 5 ticks
	for i := 0; i < 5; i++ {
		sleeper.Sleep(50 * time.Millisecond)
		time.Sleep(10 * time.Millisecond) // Give goroutine time to write
	}

	ticker.Stop()

	output := buf.String()
	dotCount := strings.Count(output, ".")
	// Should have at least a few dots (may vary due to timing)
	require.Greater(t, dotCount, 0, "Expected periodic dots in output")
	require.True(t, strings.HasSuffix(output, "\n"))
}

func TestTicker_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("", &buf)

	ticker.Start()
	time.Sleep(10 * time.Millisecond)
	ticker.Stop()

	// Even with empty message, should handle gracefully
	require.Contains(t, buf.String(), "\n")
}

func TestTicker_ConcurrentStartStop(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Concurrent", &buf)

	// Start and stop rapidly
	for i := 0; i < 10; i++ {
		ticker.Start()
		time.Sleep(1 * time.Millisecond)
		ticker.Stop()
	}

	// Should not panic or race
	require.False(t, ticker.active)
}
