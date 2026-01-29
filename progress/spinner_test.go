package progress

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockSleeper is a test double for Sleeper that tracks sleep calls
type MockSleeper struct {
	sleepCalls []time.Duration
}

func (m *MockSleeper) Sleep(d time.Duration) {
	m.sleepCalls = append(m.sleepCalls, d)
}

func TestDefaultSleeper_Sleep(t *testing.T) {
	sleeper := DefaultSleeper{}
	start := time.Now()
	sleeper.Sleep(10 * time.Millisecond)
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed, 10*time.Millisecond, "should sleep for at least the specified duration")
}

func TestNewSpinner(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"a", "b", "c"}
	delay := 50 * time.Millisecond
	message := "Testing"

	spinner := NewSpinner(message, chars, delay, &buf)

	assert.Equal(t, message, spinner.message)
	assert.Equal(t, chars, spinner.chars)
	assert.Equal(t, delay, spinner.delay)
	assert.Equal(t, &buf, spinner.writer)
	assert.False(t, spinner.active)
	assert.NotNil(t, spinner.stopChan)
	assert.IsType(t, DefaultSleeper{}, spinner.sleeper)
}

func TestNewSpinnerWithSleeper(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"a", "b"}
	delay := 30 * time.Millisecond
	message := "Custom"
	mockSleeper := &MockSleeper{}

	spinner := NewSpinnerWithSleeper(message, chars, delay, &buf, mockSleeper)

	assert.Equal(t, message, spinner.message)
	assert.Equal(t, chars, spinner.chars)
	assert.Equal(t, delay, spinner.delay)
	assert.Equal(t, &buf, spinner.writer)
	assert.Equal(t, mockSleeper, spinner.sleeper)
	assert.False(t, spinner.active)
}

func TestNewDefaultSpinner(t *testing.T) {
	message := "Loading"
	spinner := NewDefaultSpinner(message)

	assert.Equal(t, message, spinner.message)
	assert.Len(t, spinner.chars, 8, "should have 8 animation characters")
	assert.Equal(t, 100*time.Millisecond, spinner.delay)
	assert.False(t, spinner.active)
}

func TestNewDefaultSpinnerWithOutput(t *testing.T) {
	var buf bytes.Buffer
	message := "Processing"
	spinner := NewDefaultSpinnerWithOutput(message, &buf)

	assert.Equal(t, message, spinner.message)
	assert.Len(t, spinner.chars, 8)
	assert.Equal(t, 100*time.Millisecond, spinner.delay)
	assert.Equal(t, &buf, spinner.writer)
	assert.False(t, spinner.active)
}

func TestSpinner_StartAndStop(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"a", "b", "c"}
	delay := 10 * time.Millisecond
	mockSleeper := &MockSleeper{}

	spinner := NewSpinnerWithSleeper("Test", chars, delay, &buf, mockSleeper)

	// Start the spinner
	spinner.Start()
	assert.True(t, spinner.active, "spinner should be active after Start()")

	// Let it run for a bit
	time.Sleep(50 * time.Millisecond)

	// Stop the spinner
	spinner.Stop()
	assert.False(t, spinner.active, "spinner should be inactive after Stop()")

	// Verify sleeper was called
	assert.NotEmpty(t, mockSleeper.sleepCalls, "sleeper should have been called")

	// Verify all sleep calls were with the correct delay
	for _, d := range mockSleeper.sleepCalls {
		assert.Equal(t, delay, d, "all sleep calls should use the configured delay")
	}
}

func TestSpinner_StartTwice(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Test", []string{"a"}, 10*time.Millisecond, &buf, mockSleeper)

	spinner.Start()
	assert.True(t, spinner.active)

	// Starting again should do nothing
	spinner.Start()
	assert.True(t, spinner.active)

	spinner.Stop()
}

func TestSpinner_StopWhenNotActive(t *testing.T) {
	var buf bytes.Buffer
	spinner := NewSpinner("Test", []string{"a"}, 10*time.Millisecond, &buf)

	// Stopping an inactive spinner should not panic
	assert.NotPanics(t, func() {
		spinner.Stop()
	})
	assert.False(t, spinner.active)
}

func TestSpinner_Animation(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"1", "2", "3"}
	mockSleeper := &MockSleeper{}
	message := "Loading"

	spinner := NewSpinnerWithSleeper(message, chars, 10*time.Millisecond, &buf, mockSleeper)

	spinner.Start()
	time.Sleep(50 * time.Millisecond)
	spinner.Stop()

	output := buf.String()

	// Should contain the message
	assert.Contains(t, output, message, "output should contain the message")

	// Should have cycled through characters multiple times
	assert.GreaterOrEqual(t, len(mockSleeper.sleepCalls), 3, "should have completed at least one full cycle")
}

func TestSpinner_Erase(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"a", "b"}
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Test", chars, 5*time.Millisecond, &buf, mockSleeper)

	spinner.Start()
	time.Sleep(30 * time.Millisecond)
	spinner.Stop()

	output := buf.String()

	// Output should contain backspace characters for erasing
	assert.Contains(t, output, "\b", "output should contain backspace characters")
}

func TestSpinner_MultibyteCharacters(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"⣾", "⣽", "⣻"}
	mockSleeper := &MockSleeper{}
	message := "Processing"

	spinner := NewSpinnerWithSleeper(message, chars, 5*time.Millisecond, &buf, mockSleeper)

	spinner.Start()
	time.Sleep(30 * time.Millisecond)
	spinner.Stop()

	output := buf.String()

	// Should handle multibyte characters correctly
	assert.NotEmpty(t, output)
	// The erase function should work with multibyte characters
	assert.Contains(t, output, "\b", "should erase multibyte characters")
}

func TestSpinner_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	chars := []string{"a", "b"}
	mockSleeper := &MockSleeper{}

	spinner := NewSpinnerWithSleeper("", chars, 5*time.Millisecond, &buf, mockSleeper)

	spinner.Start()
	time.Sleep(20 * time.Millisecond)
	spinner.Stop()

	// Should not panic with empty message
	assert.NotPanics(t, func() {
		spinner.Start()
		time.Sleep(20 * time.Millisecond)
		spinner.Stop()
	})
}

func TestSpinner_ConcurrentStartStop(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Test", []string{"a", "b"}, 5*time.Millisecond, &buf, mockSleeper)

	// Multiple concurrent start/stop should be safe
	for i := 0; i < 3; i++ {
		spinner.Start()
		time.Sleep(10 * time.Millisecond)
		spinner.Stop()
		time.Sleep(5 * time.Millisecond)
	}

	assert.False(t, spinner.active)
}
