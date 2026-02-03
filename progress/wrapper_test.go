package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWrapper(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Test", &buf)
	wrapper := newWrapper(ticker)

	assert.Equal(t, ticker, wrapper.ticker)
}

func TestNewDefaultWrapper(t *testing.T) {
	message := "Loading"
	wrapper := NewDefaultWrapper(message)

	// Message should have "..." appended
	assert.Equal(t, "Loading...", wrapper.ticker.message)
	assert.Equal(t, 5*time.Second, wrapper.ticker.interval)
}

func TestNewDefaultWrapper_MessageWithPeriod(t *testing.T) {
	message := "Loading."
	wrapper := NewDefaultWrapper(message)

	// Should not add ellipsis if message already ends with period
	assert.Equal(t, "Loading.", wrapper.ticker.message)
}

func TestNewDefaultWrapperWithOutput(t *testing.T) {
	var buf bytes.Buffer
	message := "Processing"
	wrapper := NewDefaultWrapperWithOutput(message, &buf)

	assert.Equal(t, "Processing...", wrapper.ticker.message)
	assert.Equal(t, &buf, wrapper.ticker.writer)
}

func TestWrapper_WrapAction(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	ticker := NewTickerWithSleeper("Working", 50*time.Millisecond, &buf, mockSleeper)
	wrapper := newWrapper(ticker)

	actionCalled := false
	wrapper.WrapAction(func() {
		// Simulate some work
		for i := 0; i < 3; i++ {
			mockSleeper.Sleep(50 * time.Millisecond)
		}
		actionCalled = true
	})

	assert.True(t, actionCalled, "action should have been called")
	assert.False(t, wrapper.ticker.active, "ticker should be stopped after action completes")

	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "Working"), "should start with message")
	assert.True(t, strings.HasSuffix(output, "\n"), "should end with newline")
}

func TestWrapper_WrapAction_WithPeriodicDots(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	ticker := NewTickerWithSleeper("Progress...", 25*time.Millisecond, &buf, mockSleeper)
	wrapper := newWrapper(ticker)

	wrapper.WrapAction(func() {
		// Simulate work with periodic sleeps
		for i := 0; i < 5; i++ {
			mockSleeper.Sleep(25 * time.Millisecond)
			time.Sleep(5 * time.Millisecond) // Give goroutine time to print
		}
	})

	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "Progress..."), "should start with message")
	dotCount := strings.Count(output, ".")
	// Should have dots from message (3) plus periodic dots
	assert.Greater(t, dotCount, 3, "should have periodic dots in addition to message")
	assert.True(t, strings.HasSuffix(output, "\n"), "should end with newline")
}

func TestWrapper_WrapAction_ActionPanics(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Test", &buf)
	wrapper := newWrapper(ticker)

	assert.Panics(t, func() {
		wrapper.WrapAction(func() {
			panic("test panic")
		})
	}, "should propagate panic from action")
}

func TestWrapper_WrapAction_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	wrapper := NewDefaultWrapperWithOutput("", &buf)

	wrapper.WrapAction(func() {
		time.Sleep(10 * time.Millisecond)
	})

	output := buf.String()
	// Empty message gets "..." appended, then newline
	assert.Equal(t, "...\n", output)
}

func TestWrapper_WrapAction_MultipleActions_Sequential(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultTicker("Task...", &buf)
	wrapper := newWrapper(ticker)

	count := 0
	for i := 0; i < 3; i++ {
		wrapper.WrapAction(func() {
			count++
			time.Sleep(5 * time.Millisecond)
		})
	}

	assert.Equal(t, 3, count, "should execute all actions")
	assert.False(t, wrapper.ticker.active, "ticker should be stopped after last action")

	// Should have 3 separate ticker outputs (3 messages + 3 newlines)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.GreaterOrEqual(t, len(lines), 3, "should have output from all actions")
}
